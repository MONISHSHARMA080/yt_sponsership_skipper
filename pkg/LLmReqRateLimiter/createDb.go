package llmreqratelimiter

import (
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/go-libsql"
)

func GetRateLimiterDb(dbFilePath string) (*sql.DB, error) {
	createDBAndTableSql := `
    CREATE TABLE IF NOT EXISTS rate_limit_user(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_email TEXT NOT NULL,
    request_timestamp DATETIME NOT NULL
  );
  CREATE INDEX IF NOT EXISTS idx_user_email ON rate_limit_user(user_email);
  CREATE INDEX IF NOT EXISTS idx_request_timestamp ON rate_limit_user(request_timestamp);
  -- A composite index is often very beneficial for queries filtering by user and time
  CREATE INDEX IF NOT EXISTS idx_user_time ON rate_limit_user(user_email, request_timestamp);
  `
	println("in the DB create func")

	db, err := sql.Open("libsql", dbFilePath)
	if err != nil {
		// Return the error if opening the database fails
		return nil, fmt.Errorf("failed to open database file %s: %w", dbFilePath, err)
	}

	_, err = db.Exec(createDBAndTableSql)
	if err != nil {
		db.Close() // Close the connection if setup fails
		return nil, fmt.Errorf("failed to set up database schema: %w", err)
	}

	// Return the database connection and no error
	return db, nil
}
