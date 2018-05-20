package todo

import "github.com/eyecuelab/kit/set"

const (
	STATUS_IN_PROGRESS = "In Progress"
	STATUS_NEW         = "New"
	STATUS_CLOSED      = "Closed"
)

var allowedStatuses = set.FromStrings(STATUS_NEW, STATUS_IN_PROGRESS, STATUS_CLOSED)

// Todos is a list of todo.Todo structs
type Todos struct {
	TodoList []Todo `json:"todos"`
}

// Todo is a struct containing the ID of a todo, as well as, title and status
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// CreateOrUpdateTodo is the expected payload for a create or update todo request
type CreateOrUpdateTodo struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}
