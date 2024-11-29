package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/mikhatanu/dating-test/db"
	rest_api "github.com/mikhatanu/dating-test/rest-api"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const (
	createUserTable = `
		CREATE TABLE IF NOT EXISTS user (
			username text unique not null,
			password nvarchar(255) not null
		);
	`

	host_port = "0.0.0.0:3000"
)

func main() {
	// check for sqlite file. Create if not exist
	_, err := os.Stat("./db/app.db")
	if err != nil {
		_, err2 := os.Create("./db/app.db")
		if err2 != nil {
			log.Fatal(err2.Error())
		}
	}

	// Open sql connection, save to global variable in db package
	db.DB, db.DB_error = sql.Open("sqlite3", "./db/app.db")
	if db.DB_error != nil {
		log.Fatal(db.DB_error)
	}
	defer db.DB.Close()

	// Initialize the db with user table
	_, errExec := db.DB.Exec(createUserTable)
	if errExec != nil {
		log.Println(errExec.Error())
		log.Fatalf("Error creating user table")
	}

	// create server mux and set the handle
	mux := http.NewServeMux()
	mux.Handle("/rest/v1/login", http.HandlerFunc(rest_api.Login))
	mux.Handle("/rest/v1/signup", http.HandlerFunc(rest_api.Signup))

	// start the server
	log.Printf("listening on port %v", host_port)
	http.ListenAndServe(host_port, mux)
}
