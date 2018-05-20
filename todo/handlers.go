package todo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Create will allow a user to create a new todo
// The supported body is {"title": "", "status": ""}
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if resp, err := create(openDB(), r.Body); err != nil {
		writeErr(w, err)
		log.Print(err)
	} else {
		writeOKResp(w, resp)
		log.Print(resp)
	}
}

func (t *CreateOrUpdateTodo) validate() error {
	if t.Status == "" || t.Title == "" {
		return BadRequestError("Todo request missing status or title")
	}

	if !allowedStatuses.Contains(t.Status) {
		return BadRequestError(fmt.Sprintf("the status %s is not supported", t.Status))
	}
	return nil
}

// Internal implementation of Create endpoint.
func create(db *sql.DB, body io.Reader) (resp []byte, err error) {
	var todo CreateOrUpdateTodo

	if err = json.NewDecoder(body).Decode(&todo); err != nil {
		return nil, BadRequestError("improperly formatted http request for 'create todo'")
	}
	if err := todo.validate(); err != nil {
		return nil, err
	}

	insertStmt := fmt.Sprintf(`INSERT INTO todo (title, status) VALUES ('%s', '%s') RETURNING id`, todo.Title, todo.Status)

	var todoID int

	// Insert and get back newly created todo ID
	if err := db.QueryRow(insertStmt).Scan(&todoID); err != nil {
		return nil, fmt.Errorf("Failed to save to db: %v", err)
	}

	fmt.Printf("Todo Created -- ID: %d\n", todoID)

	newTodo := Todo{}
	db.QueryRow("SELECT id, title, status FROM todo WHERE id=$1", todoID).Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)

	return json.Marshal(newTodo)
}

// List will provide a list of all current to-dos
func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if resp, err := list(openDB()); err != nil {
		writeErr(w, err)
		log.Print(err)
	} else {
		writeOKResp(w, resp)
		log.Print(resp)
	}
}

//Internal implementation of List endpoint
func list(db *sql.DB) (resp []byte, err error) {
	todoList := []Todo{}

	rows, err := db.Query("SELECT id, title, status FROM todo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		todo := Todo{}
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
			return nil, err
		}
		todoList = append(todoList, todo)
	}

	return json.Marshal(Todos{TodoList: todoList})
}

func update(db *sql.DB, body io.Reader, id int) (resp []byte, err error) {
	var todo CreateOrUpdateTodo

	if err = json.NewDecoder(body).Decode(&todo); err != nil {
		return nil, BadRequestError("improperly formatted http request for 'update todo'")
	} else if err := todo.validate(); err != nil {
		return nil, err
	}

	const updateStatement = `UPDATE todo SET title = $2, status = $3 WHERE id = $1;`
	if _, err := db.Exec(`UPDATE todo SET title = $2, status = $3 WHERE id = $1;`, id, todo.Title, todo.Status); err != nil {
		return nil, err
	}

	fmt.Printf("Todo Updated -- ID: %d\n", id)

	var updated Todo
	db.QueryRow("SELECT id, title, status FROM todo WHERE id=$1", id).Scan(&updated.ID, &updated.Title, &updated.Status)

	return json.Marshal(updated)
}
func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
