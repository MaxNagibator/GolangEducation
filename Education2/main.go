package main

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand/v2"

	_ "github.com/lib/pq"
)

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

func main() {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "RjirfLeyz"
		dbname   = "money-dev2"
	)

	databaseConnectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var dbProvider = sqlProvider{
		connectionString: databaseConnectionString,
		state:            false,
	}

	dbProvider.ExecuteNonQuery("DELETE FROM public.operations WHERE user_id = $1", 3)

	nextOperatonId := dbProvider.QueryInt("SELECT next_operation_id FROM public.domain_users WHERE id = 3")
	nextOperatonId++

	insertedSum := toFixed(rand.Float64()*1000, 2)
	insertedComment := "test"
	insertedCount := dbProvider.ExecuteNonQuery("INSERT INTO public.operations (user_id, id, sum, comment, category_id, date, is_deleted) VALUES ($1, $2, $3, $4, $5, NOW(), false)", 3, nextOperatonId, insertedSum, insertedComment, 1)
	fmt.Println("inserted ", insertedCount)
	updatedCount := dbProvider.ExecuteNonQuery("UPDATE public.domain_users SET next_operation_id = $1 WHERE id = 3", nextOperatonId)
	fmt.Println("updated ", updatedCount)

	rows, err := dbProvider.ExecuteQuery("SELECT id, sum, comment FROM public.operations WHERE user_Id = 3")
	if err != nil {
		fmt.Println("Error execute: %v\n", err)
		return
	}
	rowIndex := 0
	for rows.Next() {
		rowIndex++
		var id int
		var sum float32
		var comment string
		rows.Scan(&id, &sum, &comment)
		fmt.Println(rowIndex, id, sum, comment)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
