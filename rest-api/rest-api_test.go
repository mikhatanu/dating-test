package rest_api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mikhatanu/dating-test/db"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func TestREST(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/v1/login", http.HandlerFunc(Login))
	mux.HandleFunc("/rest/v1/signup", http.HandlerFunc(Signup))
	db.DB, db.DB_error = sql.Open("sqlite3", "./app.db")
	if db.DB_error != nil {
		log.Fatal(db.DB_error)
	}
	defer db.DB.Close()
	defer os.Remove("./app.db")

	createUserTable := `
		CREATE TABLE IF NOT EXISTS user (
			username text unique not null,
			password nvarchar(255) not null
		);
	`
	_, err := db.DB.Exec(createUserTable)
	if err != nil {
		log.Fatalln(err.Error())
	}

	signupData := []struct {
		testname string
		username string
		password string
		want     int
	}{
		{
			testname: "signup ok",
			username: "test",
			password: "test",
			want:     http.StatusOK,
		},
		{
			testname: "signup empty username",
			username: "",
			password: "test",
			want:     http.StatusBadRequest,
		},
		{
			testname: "signup empty password",
			username: "test",
			password: "",
			want:     http.StatusBadRequest,
		},
	}
	for _, tt := range signupData {
		t.Run(tt.testname, func(t *testing.T) {
			postData := map[string]interface{}{
				"username": tt.username,
				"password": tt.password,
			}
			resp := httptest.NewRecorder()
			body, _ := json.Marshal(postData)
			req := httptest.NewRequest(http.MethodPost, "/rest/v1/signup", bytes.NewReader(body))

			mux.ServeHTTP(resp, req)

			got := resp.Result().StatusCode
			want := tt.want
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}

	loginData := []struct {
		testname string
		username string
		password string
		want     int
	}{
		{
			testname: "correct login",
			username: "test",
			password: "test",
			want:     http.StatusOK,
		},
		{
			testname: "wrong login",
			username: "error",
			password: "test",
			want:     http.StatusBadRequest,
		},
		{
			testname: "empty username",
			username: "",
			password: "test",
			want:     http.StatusBadRequest,
		},
		{
			testname: "empty password",
			username: "test",
			password: "",
			want:     http.StatusBadRequest,
		},
	}
	for _, tt := range loginData {
		t.Run(tt.testname, func(t *testing.T) {
			resp := httptest.NewRecorder()
			post := map[string]interface{}{
				"username": tt.username,
				"password": tt.password,
			}
			body, _ := json.Marshal(post)
			req := httptest.NewRequest(http.MethodPost, "/rest/v1/login", bytes.NewReader(body))

			mux.ServeHTTP(resp, req)

			got := resp.Result().StatusCode
			want := tt.want
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
