package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

type User struct {
	ID   int
	Name string
}

func queryUsers(db *sql.DB) []User {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query users %s", err)
		os.Exit(1)
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan user %s", err)
			return nil
		}

		users = append(users, user)
		fmt.Fprintf(os.Stdout, "user: %s\n", user.Name)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to scan user %s", err)
	}

	return users
}

func main() {
	dbName := "file:./local.db"

	db, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}
	defer db.Close()

	users := queryUsers(db)

	for _, user := range users {
		fmt.Fprintf(os.Stdout, "user: %s\n", user.Name)
	}
}
