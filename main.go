package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var db *sqlx.DB

type person struct {
	Id   int64
	Blog string
}

var port string = "4000"

func main() {
	var err error
	connStr := "postgres://root:GzWPeJv6HODIqI1mjH8UyvvprIGc3rAv@dpg-co6irf8l6cac73a9p1p0-a.singapore-postgres.render.com:5432/blogging_kfua"
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	people := []person{}
	db.Select(&people, "SELECT * FROM allblogs ")
	for _, all := range people {
		fmt.Printf("users id: %d ,name: %s\n", all.Id, all.Blog)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is backend code for crud"))
	})
	r.Route("/blogs", func(r chi.Router) {
		r.Get("/", getUsersAll)
		r.Post("/", createUser)
		r.Put("/{id}", updateUser)
		r.Delete("/{id}", deleteUser)
	})

	log.Printf("Listening on port %s", port)
	http.ListenAndServe(":"+port, r)

}

func getUsersAll(w http.ResponseWriter, r *http.Request) {
	var people []person

	query := "SELECT *  FROM allblogs"
	err := db.Select(&people, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(people)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var people person
	err := json.NewDecoder(r.Body).Decode(&people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	query := "INSERT INTO allblogs (blog) VALUES($1) RETURNING  id"
	err = db.QueryRow(query, people.Blog).Scan(&people.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(people)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := "DELETE FROM allblogs WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("task with ID  %s  deleted successfully", id)))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var people person
	err := json.NewDecoder(r.Body).Decode(&people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	query := "UPDATE allblogs SET blog = $1 WHERE id = $2 RETURNING *"
	_, err = db.Exec(query, people.Blog, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("task with ID  %s  Updated successfully", id)))
}
