package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Database connection string (replace with your credentials)
const dsn = "root:Amaan@123@tcp(127.0.0.1:3306)/toronto_time_db"

// Struct to represent the JSON response
type TimeResponse struct {
	CurrentTime string `json:"current_time"`
}

func main() {
	// Connect to the MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Create /current-time endpoint
	http.HandleFunc("/current-time", func(w http.ResponseWriter, r *http.Request) {
		// Get the current time in Toronto timezone
		loc, err := time.LoadLocation("America/Toronto")
		if err != nil {
			http.Error(w, "Failed to load timezone", http.StatusInternalServerError)
			log.Printf("Timezone error: %v", err)
			return
		}
		torontoTime := time.Now().In(loc)

		// Log the time to the database
		_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", torontoTime)
		if err != nil {
			http.Error(w, "Failed to log time to database", http.StatusInternalServerError)
			log.Printf("Database error: %v", err)
			return
		}

		// Create the JSON response
		response := TimeResponse{
			CurrentTime: torontoTime.Format("2006-01-02 15:04:05"),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Start the server
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
