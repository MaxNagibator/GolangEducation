package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	_ "github.com/lib/pq"
)

func main() {

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/debts", &DebtsHandler{})
	mux.Handle("/debts/", &DebtsHandler{})

	// Run the server
	http.ListenAndServe(":4545", mux)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

var (
	DebtRe       = regexp.MustCompile(`^/debts/*$`)
	DebtReWithID = regexp.MustCompile(`^/debts/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

type debtService struct {
	sqlProvider sqlProvider
}

func (ds *debtService) Add(debt Debt) error {
	_, err := ds.sqlProvider.ExecuteQuery("INSERT INTO public.debts(id, name, status) VALUES($1, $2, $3)", debt.Id, debt.Name, debt.Status)
	return err
}

type DebtsHandler struct {
	debtService debtService
}

func (h *DebtsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && DebtRe.MatchString(r.URL.Path):
		h.CreateDebt(w, r)
		return
	case r.Method == http.MethodGet && DebtRe.MatchString(r.URL.Path):
		h.ListDebts(w, r)
		return
	case r.Method == http.MethodGet && DebtReWithID.MatchString(r.URL.Path):
		h.GetDebt(w, r)
		return
	case r.Method == http.MethodPut && DebtReWithID.MatchString(r.URL.Path):
		h.UpdateDebt(w, r)
		return
	case r.Method == http.MethodDelete && DebtReWithID.MatchString(r.URL.Path):
		h.DeleteDebt(w, r)
		return
	default:
		return
	}
}

func (h *DebtsHandler) CreateDebt(w http.ResponseWriter, r *http.Request) {
	var debt Debt
	if err := json.NewDecoder(r.Body).Decode(&debt); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	//resourceID := slug.Make(recipe.Name)
	if err := h.debtService.Add(debt); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h *DebtsHandler) ListDebts(w http.ResponseWriter, r *http.Request)  {}
func (h *DebtsHandler) GetDebt(w http.ResponseWriter, r *http.Request)    {}
func (h *DebtsHandler) UpdateDebt(w http.ResponseWriter, r *http.Request) {}
func (h *DebtsHandler) DeleteDebt(w http.ResponseWriter, r *http.Request) {}

//CREATE TABLE IF NOT EXISTS public.debts
//(
//    id integer NOT NULL,
//    comment character varying(4000) COLLATE pg_catalog."default" NOT NULL,
//    status integer NOT NULL,
//    CONSTRAINT pk_debts PRIMARY KEY (id)
//)

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

type Debt struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Status int    `json:"status"`
}

type sqlProvider struct {
	connectionString string
	state            bool
	db               *sql.DB
}

func (provider sqlProvider) QueryInt(query string, args ...any) int64 {
	provider.OpenConnection()
	sqlRow := provider.db.QueryRow(query, args)
	var val int64
	sqlRow.Scan(&val)
	return val
}

func (provider sqlProvider) ExecuteNonQuery(query string, args ...any) int64 {
	provider.OpenConnection()
	result, err := provider.db.Exec(query, args...)
	if err != nil {
		fmt.Println("Error execute: %v\n", err)
		return -1
	}
	r, _ := result.RowsAffected()
	return r
}

func (provider sqlProvider) ExecuteQuery(query string, args ...any) (*sql.Rows, error) {
	provider.OpenConnection()
	return provider.db.Query(query, args...)
}

func (provider *sqlProvider) OpenConnection() {
	if provider.state == false {
		db, err := sql.Open("postgres", provider.connectionString)
		if err != nil {
			fmt.Println("Unable to connect to database: %v\n", err)
			return
		}
		provider.db = db
	}
}
