package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

	query := `INSERT INTO submissions (email, apps, is_developer) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, submission.Email, submission.Apps, submission.IsDeveloper)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println("Error inserting data:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Form submission successful!"))
}
