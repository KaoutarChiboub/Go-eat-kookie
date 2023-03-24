package main

import (
	"database/sql"  // To connect to sql db
	"encoding/json" // To encode data to json
	"fmt"           // To print out information about servers and connections
	"log"           // To log out errors
	"net/http"      // To handle http requests
	"os"            // To get variables

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Machine struct {
	ID           string `json:"id"` // json mapping to be able to encode and decode json with postamn
	Name         string `json:"name"`
	OutletNumber int    `json:"outlet_number"`
	Status       string `json:"status"`
}

//type Plugs struct {
//	Number string `json:"number"`
//	pdu  string `json:"pdu"`
//}

var db *sql.DB

// The following functions are used to handle HTTP requests

func getMachines(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM machines")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		machines := []Machine{}
		for rows.Next() {
			var m Machine
			if err := rows.Scan(&m.ID, &m.Name, &m.OutletNumber, &m.Status); err != nil {
				log.Fatal(err)
			}
			if err := rows.Err(); err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(machines)
		}
	}
}

func deleteMachine(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r) //params extracts the URL parameters from the request made by the client => Here its the machine ID.
		id := params["id"]
		var m Machine
		err := db.QueryRow("SELECT * FROM machines WHERE id = ?", id).Scan(&m.ID, &m.Name, &m.OutletNumber, &m.Status)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM machines WHERE id = $1", id)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode("Machine deleted")
		}
	}
}

func getMachine(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		var m Machine
		err := db.QueryRow("SELECT * FROM machines WHERE id = $1", id).Scan(&m.ID, &m.Name, &m.OutletNumber, &m.Status)
		if err != nil {
			w.WriteHeader(http.StatusNotFound) //To set the response status code to 404 and send a header with that status code to the client
		}
		json.NewEncoder(w).Encode(m)
	}
}

func createMachine(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m Machine
		json.NewDecoder(r.Body).Decode(&m)                                                                                               //Decode the JSON body of the new machine that is sent to Postman. The decoded value will be stored in  the var <m>
		_, err := db.Exec("INSERT INTO machines (id, name, outlet, status) VALUES (?, ?, ?, ?)", m.ID, m.Name, m.OutletNumber, m.Status) //add the new machine to our table with its given characteristics
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(m)
	}
}

func updateMachine(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m Machine
		json.NewDecoder(r.Body).Decode(&m)
		params := mux.Vars(r)
		id := params["id"]
		_, err := db.Exec("UPDATE machines SET name = $1, outlet = $2, status = $3 WHERE id = $4", m.Name, m.OutletNumber, m.Status, id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(m)

	}
}

// Indide the main function we will have to define the router and functions (to do CRUD operations)

func main() {

	//Connection to postgres database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	//NewRouter is a function inside mux library to define our new router
	r := mux.NewRouter()

	//Creating our slices of movies (simplified data model for this example)

	//We call the set of functions on our 2 endpoints to handle HTTP methods (CRUD requests)
	r.HandleFunc("/users", getMachines(db)).Methods("GET")
	r.HandleFunc("/users/{id}", getMachine(db)).Methods("GET")
	r.HandleFunc("/users", createMachine(db)).Methods("POST")
	r.HandleFunc("/users/{id}", updateMachine(db)).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteMachine(db)).Methods("DELETE")

	fmt.Printf("Starting the server at port 8008\n")

	//We will start the HTTP server and set up the server to handle incoming requests using our router.
	//First, we will listen to incoming requests so that when received, it will pass it to the router to handle it.
	//With log function, we will log eventual occuring errors and exit the program. Otherwise, the execution of the program will continue until terminated.
	log.Fatal(http.ListenAndServe(":8008", r))

}
