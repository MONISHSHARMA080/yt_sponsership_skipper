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
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to open db %s", err)
    os.Exit(1)
  }
  defer db.Close()
}
