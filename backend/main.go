package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
	router.HandleFunc("api/go/users/{id}", getUser(db)).Methods("GET")
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
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
				log.Fatal(err)
			}

			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.Id, u.Name, &u.Email)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)
	}
}

func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.Id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		vars := mux.Vars(r)
		id := vars["id"]

		//execute the update query.

		_, err := db.Exec("UPDATE users SET name = $1, email = $2 HWERE id = $3", u.Name, u.Email)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch updated user from the database now.

		var updatedUser User
		err = db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&updatedUser.Id, &updatedUser.Email, &updatedUser.Name)
		if err != nil {
			log.Fatal(err)
		}

		// Show the updated user data to client.

		json.NewEncoder(w).Encode(updatedUser)
	}
}

func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.Id, &u.Name, &u.Email)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM users HERE id = $1", id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode("User deleted")
		}

	}
}
