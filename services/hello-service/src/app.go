package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const defaultVersion = "v0.0.0"

var db *sql.DB

func getVersion() string {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		return defaultVersion
	}
	return version
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":3000"
	}
	return ":" + port
}

func initDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Create a simple table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS greetings (
			id SERIAL PRIMARY KEY,
			message TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		INSERT INTO greetings (message) VALUES 
		('Hello World from K3s!');
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	var message string
	err := db.QueryRow("SELECT message FROM greetings ORDER BY created_at DESC LIMIT 1").Scan(&message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s\n", message)
}

// just example of code, not prod version
func main() {
	version := getVersion()
	log.Printf("Starting hello-service version: %s", version)

	initDB()
	defer db.Close()

	http.HandleFunc("/", helloHandler)

	port := getPort()
	fmt.Printf("Server running at http://0.0.0.0%s/\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
