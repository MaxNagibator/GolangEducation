package main

import (
	"database/sql"
	"fmt"
	"math"

	"math/rand"

	_ "github.com/lib/pq"
)

func main() {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "RjirfLeyz"
		dbname   = "money-dev2"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Unable to connect to database: %v\n", err)
		return
	}

	nextOperationIdRow := db.QueryRow("SELECT next_operation_id FROM public.domain_users WHERE id = 3")
	var nextOperatonId int
	nextOperationIdRow.Scan(&nextOperatonId)
	nextOperatonId++

	insertedSum := toFixed(rand.Float64()*1000, 2)
	insertedComment := "test"
	db.Exec("DELETE FROM public.operations WHERE user_id = 3")
	result, err := db.Exec("INSERT INTO public.operations (user_id, id, sum, comment, category_id, date, is_deleted) VALUES ($1, $2, $3, $4, $5, NOW(), false)", 3, nextOperatonId, insertedSum, insertedComment, 1)
	if err != nil {
		fmt.Println("Error execute: %v\n", err)
		return
	}
	rows22, _ := result.RowsAffected()
	if rows22 > 0 {
		fmt.Println("success", rows22)
		result, err := db.Exec("UPDATE public.domain_users SET next_operation_id = $1 WHERE id = 3", nextOperatonId)
		if err != nil {
			fmt.Println("Error execute: %v\n", err)
			return
		}
		rows22, _ := result.RowsAffected()
		fmt.Println("success2 ", rows22)
	}

	rows, err := db.Query("SELECT id, sum, comment FROM public.operations WHERE user_Id = 3")
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
