package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
  "database/sql"
  _ "github.com/lib/pq"
	"github.com/unrolled/render"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func init() {
  var err error
  db, err = sql.Open("postgres", "postgres://localhost/volunteer_be?sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }

  setupDB()
}

func setupDB() {
  db.Exec(`CREATE TABLE users (
    id SERIAL,
    user_name VARCHAR(60),
    user_email VARCHAR(60),
    user_password VARCHAR(60),
    creaet_at TIMESTAMP WITH TIME ZONE,
    user_last_login TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (id),
    CONSTRAINT users_email UNIQUE (user_email)
    );`)

  db.Exec(`INSERT INTO users (user_name, user_email, user_password) VALUES (
    'john', 'john@example.com', 'supersecret'
    );`)
}

func main() {
  defer db.Close()

  mux := http.NewServeMux()
  n := negroni.Classic()

  store := cookiestore.New([]byte("ohhhsoosecret"))
  n.Use(sessions.Sessions("global_session_store", store))

  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
      SimplePage(w, r, "login")
    } else if r.Method == "POST" {
      SignupPost(w, r)
    }
  })

  mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
      SimplePage(w, r, "login")
    } else if r.Method == "POST" {
      LoginPost(w, r)
    }
  })

  mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
      SimplePage(w, r, "signup")
    } else if r.Method == "POST" {
      SignupPost(w, r)
    }
  })

  mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    Logout(w, r)
  })

  mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
    SimpleAuthenticatedPage(w, r, "home")
  })

  mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
    APIHandler(w, r)
  })

  mux.HandleFunc("/static/", http.StripPrefix("/stati/", http.FileServer(http.Dir("static"))))

  n.UseHandler(mux)
  port := os.Getenv("PORT")
  if port == "" {
    port = "3000"
  }
  n.Run(":" + port)
}
