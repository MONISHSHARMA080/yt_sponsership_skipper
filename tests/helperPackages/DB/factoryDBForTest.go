package DB

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// this function will create a new sqlite Db in the file system and will set the env to the local DB
func CreateDBForTest() error {
	DBFileName := "localTest.db"
	println("the env for the sqlite file is ->", os.Getenv("TURSO_DATABASE_URL"), " and we in the testing environment ->", os.Getenv("IS_THIS_TESTING_ENVIRONMENT"))
	os.Setenv("TURSO_DATABASE_URL", "file:../"+DBFileName)
	os.Setenv("IS_THIS_TESTING_ENVIRONMENT", "true")
	println("the env for the sqlite file is ->", os.Getenv("TURSO_DATABASE_URL"), " and we in the testing environment ->", os.Getenv("IS_THIS_TESTING_ENVIRONMENT"))
	err := checkIfDBIsThereInTheRootDirIfNotCreateIt(DBFileName)
	if err != nil {
		return err
	}
	err = verifyThatTheTablesOnDBIsCorrect()
	if err != nil {
		println("error in checking that the DB tables is correct is ->", err.Error())
		return err
	}
	// now create the local sqlite and cp it to the root dir(not the test one as I want the go server to use it too)
	return nil
}

// func will go in the root dir and will check if the sqlite file is present if it is not then we will create it and you can simple use the DBConnect method,
// (assuming you have set the env to correct) it will get get the DB name form that
func checkIfDBIsThereInTheRootDirIfNotCreateIt(DBFileName string) error {
	dbPathToRootDir := "../" // cause when we execute it we are in the test dir so leave the one ../ out

	dirContent, err := os.ReadDir(dbPathToRootDir)
	if err != nil {
		return err
	}
	// var sqliteDBFileForTest os.DirEntry
	for i, content := range dirContent {
		if !content.IsDir() && strings.Contains(content.Name(), DBFileName) {
			println("(found the file at and checking the tables  on it)the content at : ", i, " is a file and called ", content.Name())
			println("the DB is found so will not create a new one")
			return nil
		}
	}
	// well the file is not there in the DB create one
	println("sqlite file for test not found so we are creating one ")

	newDBFile, err := os.Create("../" + DBFileName)
	if err != nil {
		return err
	}
	println("the new DB File name is ->", newDBFile.Name())
	return nil
}

// if the tables on the DB does not exist then we will create them
func verifyThatTheTablesOnDBIsCorrect() error {
	db := DbConnect()
	defer db.Close()
	tables := []struct {
		name       string
		createStmt string
	}{
		{
			name: "UserAccount",
			createStmt: `CREATE TABLE IF NOT EXISTS UserAccount (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				accountid BIGINT NOT NULL,
				email TEXT NOT NULL UNIQUE,
				userName TEXT NOT NULL,
				is_a_paid_user BOOLEAN NOT NULL DEFAULT FALSE
			)`,
		},
		{
			name: "temporaryFieldToVerifyParymentLater",
			createStmt: `CREATE TABLE IF NOT EXISTS temporaryFieldToVerifyParymentLater (
				user_account_id INT PRIMARY KEY,
				recurring_order_id TEXT NOT NULL,
				onetime_order_id TEXT NOT NULL,
				FOREIGN KEY (user_account_id) REFERENCES UserAccount(id)
			)`,
		},
		{
			name: "messageForTheUserAfterPaymentCaptured",
			createStmt: `CREATE TABLE IF NOT EXISTS messageForTheUserAfterPaymentCaptured (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_tier TEXT NOT NULL CHECK (
					user_tier IN ('recurring', 'free tier', 'one time')
				),
				razorpay_payment_id TEXT NOT NULL,
				check_for_key_update_on TIMESTAMP NOT NULL,
				version INTEGER NOT NULL DEFAULT 0,
				user_account_id INT NOT NULL,
				FOREIGN KEY (user_account_id) REFERENCES UserAccount (id)
			)`,
		},
	}

	// Check for each table and create if not exists
	for _, table := range tables {
		// Check if the table exists
		var tableName string
		query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", table.name)
		err := db.QueryRow(query).Scan(&tableName)

		if err == sql.ErrNoRows {
			// Table doesn't exist, create it
			_, err := db.Exec(table.createStmt)
			if err != nil {
				db.Close()
				fmt.Printf("Created table: %s\n", table.name)
				return fmt.Errorf("failed to create table %s: %w", table.name, err)
			}
		} else if err != nil {
			db.Close()
			return nil
		} else {
			fmt.Printf("Table exists: %s\n", table.name)
			return nil
		}
	}
	return nil
}
