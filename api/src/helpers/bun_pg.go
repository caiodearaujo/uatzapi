package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	_ "github.com/uptrace/bun/driver/pgdriver"
	"os"
	"whatsgoingon/data"
)

func getPostgresConnection() *sql.DB {
	dbUser := os.Getenv("pg_username")
	dbPwd := os.Getenv("pg_password")
	dbTCPHost := os.Getenv("pg_hostname")
	dbPort := os.Getenv("pg_port")
	dbName := "uatzapi"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	return db
}

func GetBunConnection() *bun.DB {
	pgDB := getPostgresConnection()
	return bun.NewDB(pgDB, pgdialect.New())
}

// CreateTablesFromDataPkg Get all structs from a package and create Tables in the database
func CreateTablesFromDataPkg() {
	structs := data.TablesPostgres()

	db := GetBunConnection()

	for _, structType := range structs {
		// Create table if not exists
		_, err := db.NewCreateTable().Model(structType).IfNotExists().Exec(context.Background())
		if err != nil {
			failOnError(err, "Failed to create table")
		} else {
			fmt.Printf("Table created successfully: %s", structType)
		}

	}
}
