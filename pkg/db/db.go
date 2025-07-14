package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/tursodatabase/go-libsql"
)

type Repository struct {
	db  *sql.DB
	ctx context.Context
}

func NewRepository(ctx context.Context, dbPath string) *Repository {
	if dbPath == "" {
		dbPath = getDefaultDBPath()
	}
	
	if !strings.HasPrefix(dbPath, "file:") {
		dbPath = "file:" + dbPath
	}
	
	db, err := sql.Open("libsql", dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}

	return &Repository{db: db, ctx: ctx}
}

func getDefaultDBPath() string {
	return "./local.db"
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func GetDB(dbPath string) *sql.DB {
	if dbPath == "" {
		dbPath = getDefaultDBPath()
	}
	
	if !strings.HasPrefix(dbPath, "file:") {
		dbPath = "file:" + dbPath
	}
	
	db, err := sql.Open("libsql", dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}

	return db
}
