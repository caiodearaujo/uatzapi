package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	"whatsgoingon/data"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

var (
	dbOnce sync.Once
	bunInstance *bun.DB
)

// Get environment variable with a fallback value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Initialize PostgreSQL connection using environment variables.
func getPostgresConnection() (*sql.DB, error) {
	dbUser := getEnv("pg_username", "postgres")
	dbPwd := getEnv("pg_password", "postgres")
	dbTCPHost := getEnv("pg_hostname", "localhost")
	dbPort := getEnv("pg_port", "5432")
	dbName := getEnv("pg_dbname", "uatzapi")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	return db, nil
}

// GetBunConnection returns a singleton instance of bun.DB for database operations.
func GetBunConnection() *bun.DB {
	dbOnce.Do(func() {
		pgDB, err := getPostgresConnection()
		if err != nil {
			fmt.Errorf("Failed to connect to PostgreSQL: %v", err)
		}
		bunInstance = bun.NewDB(pgDB, pgdialect.New())
	})
	return bunInstance
}

// CreateTablesFromDataPkg creates tables from all structs in the `data` package
func CreateTablesFromDataPkg() {
	structs := data.TablesPostgres()

	db := GetBunConnection()
	ctx := context.Background()

	for _, structType := range structs {
		if _, err := db.NewCreateTable().Model(structType).IfNotExists().Exec(ctx); err != nil {
			fmt.Errorf("Failed to create table: %s", structType)
		} else {
			fmt.Printf("Table created successfully: %s", structType)
		}
	}
}
