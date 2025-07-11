package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

type Repository struct {
	db  *sql.DB
	ctx context.Context
}

func NewRepository(ctx context.Context, dbName string) *Repository {
	db, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}

	return &Repository{db: db, ctx: ctx}
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func GetDB(dbName string) *sql.DB {
	db, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}

	return db
}
