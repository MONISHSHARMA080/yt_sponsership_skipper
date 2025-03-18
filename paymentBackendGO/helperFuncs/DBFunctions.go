package helperfuncs

import (
	"database/sql"
	"fmt"
	"os"
)

func DbConnect() *sql.DB {
	// println(os.Getenv("TURSO_DATABASE_URL"),os.Getenv("TURSO_AUTH_TOKEN"))
	// url := "libsql://["+os.Getenv("TURSO_DATABASE_URL")+"].turso.io?authToken=["+os.Getenv("TURSO_AUTH_TOKEN")+"]"
	// url := os.Getenv("TURSO_DATABASE_URL")+".?authToken="+os.Getenv("TURSO_AUTH_TOKEN")
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	url := fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
	// println(url,"\n\n")
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	//   _, b := db.Query("CREATE TABLE UserAccount (    accountid INT,email TEXT,    strUserToken TEXT);")
	//   if b!=nil{
	//     panic(b.Error())
	//   }
	return db
}
