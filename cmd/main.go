package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/yeboahd24/rate-limiter/handler"
	"github.com/yeboahd24/rate-limiter/router"
)

func main() {
	// Initialize SQLite database connection
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "./ratelimiter.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	err = createTables(db)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Initialize your AuthHandler
	authHandler := handler.NewAuthHandler(db)

	// Setup auth routes
	router.SetupAuthRoutes(mux, authHandler)

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func createTables(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
