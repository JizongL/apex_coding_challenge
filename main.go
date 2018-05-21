package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/efronlicht/apex_coding_challenge/todo"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

//TODO: documentation pass to explain how to set up deployment & test enviornments

func main() {
	Main(todo.OpenDB)
}

//inject DB configuration dependency into handler
func withConfig(getDB func() *sql.DB, handle func(getDB func() *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params)) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handle(getDB, w, r, ps)
	},
	)
}

func Main(db func() *sql.DB) {
	router := httprouter.New()
	router.GET("/", Status)
	router.POST("/todos", withConfig(db, todo.Create))
	router.GET("/todos", withConfig(db, todo.List))
	router.POST("/update/:id", withConfig(db, todo.Update))
	log.Println("Starting server...")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Status Request Received")
	w.WriteHeader(200)
	fmt.Fprint(w, "OK\n")
}
