package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

func GetDB(dbName string) *sql.DB {
	db, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}

	return db
}

