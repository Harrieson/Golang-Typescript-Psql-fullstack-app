package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Global variable for database connection
func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXIST users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	//Router Using Mux;

	router := mux.NewRouter()
	router.HandleFunc("/api/go/user", getUsers(db)).Methods("GET")
	router.HandleFunc("api/go/users", createUser(db)).Methods("POST")
	router.HandleFunc("api/go/users/{id}", getUsers(db)).Methods("GET")
	router.HandleFunc("api/go/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("api/go/users/{id}", deleteUser(db)).Methods("DELETE")

	// Log the Router.
	enhanceRouter := enableCORS(jsonContentTypeMiddleware(router))

	//Start the server.
	log.Fatal(http.ListenAndServe(":8000", enhanceRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//set CORS Headers

		w.Header().Set("Access-Control-Allow-Origin", "") //Allow any Origin.
		w.Header().Set("Acess-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Acess-Control-Allow-Headers", "Content-Type, Authorization")

		//Check if a request passed CORS Preflight.

	})
}
