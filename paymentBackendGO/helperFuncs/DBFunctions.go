package helperfuncs

import (
	"database/sql"
	"fmt"
	"os"
	// _ "github.com/tursodatabase/libsql-client-go/libsql"
	// _ "github.com/tursodatabase/go-libsql"
)

func DbConnect() *sql.DB {
	// println(os.Getenv("TURSO_DATABASE_URL"),os.Getenv("TURSO_AUTH_TOKEN"))
	// url := "libsql://["+os.Getenv("TURSO_DATABASE_URL")+"].turso.io?authToken=["+os.Getenv("TURSO_AUTH_TOKEN")+"]"
	// url := os.Getenv("TURSO_DATABASE_URL")+".?authToken="+os.Getenv("TURSO_AUTH_TOKEN")

	url := ""
	dbURL := ""
	authToken := ""
	isThisTestingEnv := os.Getenv("IS_THIS_TESTING_ENVIRONMENT")
	if isThisTestingEnv == "true" {
		dbURL = os.Getenv("TURSO_DATABASE_URL")
		url = dbURL
		println("\n\n testing env and the url is ->", url, "\n\n")
	} else {
		// in any case we are in prod
		dbURL = os.Getenv("TURSO_DATABASE_URL")
		authToken = os.Getenv("TURSO_AUTH_TOKEN")
		url = fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
	}

	// help me with printing the files and the dir
	// dirEntry, err := os.ReadDir(".")
	// if err != nil {
	// 	panic(err)
	// }
	// for i, entry := range dirEntry {
	// 	println("the entry at :", i, "is ->", entry.Name())
	// }
	println("the db url is ->", dbURL, " and the url is ->", url)

	// println(url,"\n\n")
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	return db
}
