package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {

	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "RjirfLeyz"
		dbname   = "go-crud"
	)

	databaseConnectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var dbProvider = sqlProvider{
		connectionString: databaseConnectionString,
		state:            false,
	}

	var debtService = debtService{
		sqlProvider: dbProvider,
	}

	debtsHandler := DebtsHandler{
		debtService: debtService,
	}

	mux := http.NewServeMux()
	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/debts", &debtsHandler)
	mux.Handle("/debts/", &debtsHandler)

	// Run the server
	http.ListenAndServe(":4545", mux)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

var (
	DebtRe       = regexp.MustCompile(`^/debts/*$`)
	DebtReWithID = regexp.MustCompile(`^/debts/[0-9]+$`)
)

type debtService struct {
	sqlProvider sqlProvider
}

func (ds *debtService) Add(debt Debt) error {
	_, err := ds.sqlProvider.ExecuteQuery("INSERT INTO public.debts(id, name, status) VALUES($1, $2, $3)", debt.Id, debt.Name, debt.Status)
	return err
}

func (ds *debtService) Delete(debtId int) error {
	_, err := ds.sqlProvider.ExecuteQuery("DELETE FROM public.debts WHERE id = $1", debtId)
	return err
}
func (ds *debtService) GetAll() ([]Debt, error) {

	rows, err := ds.sqlProvider.ExecuteQuery("SELECT id, name, status FROM public.debts")
	if err != nil {
		fmt.Println("Error execute: %v\n", err)
		return nil, err
	}
	rowIndex := 0
	debts := []Debt{}
	for rows.Next() {
		rowIndex++
		var id int
		var name string
		var status int
		rows.Scan(&id, &name, &status)
		debt := Debt{
			Id:     id,
			Name:   name,
			Status: status,
		}
		debts = append(debts, debt)
	}

	return debts, err
}

type DebtsHandler struct {
	debtService debtService
}

func (h *DebtsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("piu pau", r.Method, r.URL.Path)
	asd := DebtReWithID.MatchString(r.URL.Path)
	fmt.Println(asd)
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
		fmt.Println(err)
		InternalServerErrorHandler(w, r)
		return
	}
	//resourceID := slug.Make(recipe.Name)
	if err := h.debtService.Add(debt); err != nil {
		fmt.Println(err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h *DebtsHandler) ListDebts(w http.ResponseWriter, r *http.Request) {
	debtList, err := h.debtService.GetAll()

	jsonBytes, err := json.Marshal(debtList)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h *DebtsHandler) GetDebt(w http.ResponseWriter, r *http.Request)    {}
func (h *DebtsHandler) UpdateDebt(w http.ResponseWriter, r *http.Request) {}
func (h *DebtsHandler) DeleteDebt(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Path[len("/debts/"):len(r.URL.Path)]
	id, err := strconv.Atoi(idString)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.debtService.Delete(id); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
