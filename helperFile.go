package main

import (
  "database/sql"
  "fmt"
  "os"
  _ "github.com/tursodatabase/go-libsql"
)

func DbConnect() {
  dbName := "file:./local.db"
  println("db --><<><>>>")

  db, err := sql.Open("libsql", dbName)
  defer db.Close()

	// Create the "user" table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, name TEXT);")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create table %s", err)
		os.Exit(1)
	}

	// Query to get all table names in the database
	res, err := db.Query("SELECT name FROM sqlite_master WHERE type = 'table';")
	if err != nil {
		panic(err.Error())
	}
	defer res.Close()

	// Iterate over the result set and print the table names
	println("Tables in the database:")
	for res.Next() {
		var tableName string
		err := res.Scan(&tableName)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(tableName)
	}

	// Check for any error after iteration
	if err = res.Err(); err != nil {
		panic(err.Error())
	}
}