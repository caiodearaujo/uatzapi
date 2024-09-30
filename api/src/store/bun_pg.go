package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	"whatsgoingon/data"
	"whatsgoingon/handler"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

var (
	dbOnce      sync.Once
	bunInstance *bun.DB
)

// Initialize PostgreSQL connection using environment variables.
func getPostgresConnection() (*sql.DB, error) {
	dbUser := os.Getenv("PG_USERNAME")
	dbPwd := os.Getenv("PG_PASSWORD")
	dbTCPHost := os.Getenv("PG_HOSTNAME")
	dbPort := os.Getenv("PG_PORT")
	dbName := os.Getenv("PG_DATABASE")
	dbSchema := os.Getenv("PG_UA_SCHEMA")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=%s&sslmode=disable", dbUser, dbPwd, dbTCPHost, dbPort, dbName, dbSchema)

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
			handler.FailOnError(err, "Failed to connect to PostgreSQL")
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
			handler.FailOnError(err, fmt.Sprintf("Failed to create table: %s", structType))
		}
	}
}
