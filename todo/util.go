package todo

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type BadRequestError string

func (err BadRequestError) Error() string {
	return string(err)
}

func writeErr(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case nil:
		panic("called writeErr on nil error") // this shouldn't happen
	case BadRequestError:
		http.Error(w, string(err), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", 500) // don't expose internal errors
	}
}

func writeOKResp(w http.ResponseWriter, jsonResp []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, jsonResp)
}

func openDB() *sql.DB {
	db, _ := sql.Open("postgres",
		fmt.Sprintf("user=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME")))
	return db
}
