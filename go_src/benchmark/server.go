package main

import (
	"fmt"
    "log"
	"net/http"
	"time"
    "encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
const DB_DSN = "root:root@/go_benchmark"

type User struct {
	ID int
    Username string
    Email string
}


func main() {
    http.HandleFunc("/json", jsonTestHandler)
    http.HandleFunc("/benchmark", benchmarkTestHandler)

	initDB();

	defer db.Close();

    fmt.Printf("Starting server at port 8081\n")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal(err)
	}
}

func jsonTestHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	user.ID = 1
	user.Username = "username"
	user.Email = "email@email.com"
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}


func benchmarkTestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			readTestHandler(w, r);
		case "POST":
			writeTestHandler(w, r);
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "not found"}`))
    }
}

func writeTestHandler(w http.ResponseWriter, r *http.Request) {
	qr, err := writeQuery()

	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qr)
}

func readTestHandler(w http.ResponseWriter, r *http.Request) {
	var user = readQuery();

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}


func initDB() {
	var err error
	db, err = sql.Open("mysql", DB_DSN)

	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
}

func readQuery() User {
	var user User

	err := db.QueryRow("SELECT id, username, email FROM users where id = ?", 1).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		log.Fatal(err.Error())
	}

	return user
}

func writeQuery() (int64, error) {
	// Without prepared statement:
	r, err := db.Exec("INSERT into users (username, email) values ('username', 'email@email.com')");

	// With prepared statement:
	//r, err := db.Prepare("INSERT into users (username, email) values (?, ?)")
	//db.Exec("username", "email@email.com")

	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}

	return r.LastInsertId()
}