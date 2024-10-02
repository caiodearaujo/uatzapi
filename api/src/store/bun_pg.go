package store

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"whatsgoingon/handler"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	dbOnce      sync.Once // Ensures that the connection is established only once (singleton).
	bunInstance *bun.DB   // Global instance of bun.DB for database operations.
)

// getPostgresConnection initializes a PostgreSQL connection using environment variables.
// Returns a *sql.DB connection or an error if the connection fails.
func getPostgresConnection() (*sql.DB, error) {
	// Fetch database configuration from environment variables.
	dbUser := os.Getenv("PG_USERNAME")
	dbPwd := os.Getenv("PG_PASSWORD")
	dbTCPHost := os.Getenv("PG_HOSTNAME")
	dbPort := os.Getenv("PG_PORT")
	dbName := os.Getenv("PG_DATABASE")
	dbSchema := os.Getenv("PG_UA_SCHEMA")

	// Create a Data Source Name (DSN) using the provided environment variables.
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=%s&sslmode=disable",
		dbUser, dbPwd, dbTCPHost, dbPort, dbName, dbSchema)

	// Open a new connection to the PostgreSQL database using Bun's pgdriver.
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Verify the connection by pinging the database.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	return db, nil
}

// GetBunConnection returns a singleton instance of *bun.DB for executing queries.
// It uses sync.Once to ensure the connection is established only once during the application's lifetime.
func GetBunConnection() *bun.DB {
	// Ensure the database connection is created only once.
	dbOnce.Do(func() {
		pgDB, err := getPostgresConnection()
		if err != nil {
			// If there is an error establishing the connection, the application will fail.
			handler.FailOnError(err, "Failed to connect to PostgreSQL")
		}
		// Initialize the Bun instance with the PostgreSQL dialect.
		bunInstance = bun.NewDB(pgDB, pgdialect.New())
	})
	return bunInstance
}
