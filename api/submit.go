package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}
}

type FormSubmission struct {
	Email       string `json:"email"`
	Apps        string `json:"apps"`
	IsDeveloper bool   `json:"isDeveloper"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var submission FormSubmission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = insertSubmission(submission)

	if err != nil {
		if strings.Contains(err.Error(), `relation "submissions" does not exist`) {
			// Table does not exist, create it and try again
			err = createTable()
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				fmt.Println("Error creating submissions table:", err)
				return
			}

			// Try inserting again after creating the table
			err = insertSubmission(submission)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				fmt.Println("Error inserting data after creating table:", err)
				return
			}
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Println("Error inserting data:", err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Form submission successful!"))
}

func insertSubmission(submission FormSubmission) error {
	query := `INSERT INTO submissions (email, apps, is_developer) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, submission.Email, submission.Apps, submission.IsDeveloper)
	return err
}

func createTable() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS submissions (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		apps TEXT NOT NULL,
		is_developer BOOLEAN NOT NULL,
		submitted_at TIMESTAMPTZ DEFAULT NOW()
	);
	`
	_, err := db.Exec(createTableQuery)
	return err
}
