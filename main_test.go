package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	challenge "github.com/efronlicht/apex_coding_challenge"
	"github.com/efronlicht/apex_coding_challenge/todo"

	"github.com/stretchr/testify/assert"
)

func Test_Main(t *testing.T) {
	db := todo.OpenTestDB()
	db.Exec(todo.DROP_TABLE)
	defer db.Exec(todo.DROP_TABLE)
	db.Exec(todo.RECREATE_SCHEMA)

	go challenge.Main(todo.OpenTestDB)
	time.Sleep(100 * time.Millisecond)

	//create some item on a blank database: we know it's item #1
	toCreate := todo.CreateOrUpdateTodo{"integration test", todo.STATUS_IN_PROGRESS}
	body, _ := json.Marshal(toCreate)
	_, err := http.Post("http://localhost:8080/todos", "application/json", bytes.NewReader(body))

	assert.NoError(t, err)

	//test that we properly send errors w/ correct content headers
	badReq := todo.CreateOrUpdateTodo{"bad body", "invalid status"}
	body, _ = json.Marshal(badReq)

	resp, _ := http.Post("http://localhost:8080/todos", "application/json", bytes.NewReader(body))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	//update our todo:
	toUpdate := todo.CreateOrUpdateTodo{"integration test", todo.STATUS_CLOSED}
	body, _ = json.Marshal(toUpdate)
	resp, _ = http.Post("http://localhost:8080/update/1", "application/json", bytes.NewReader(body))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	//now get a list of the statuses and see if it's registered our update
	var got []byte

	resp, _ = http.Get("http://localhost:8080/todos")
	resp.Body.Read(got)

	assert.Contains(t, "integration test", string(got))
	assert.Contains(t, todo.STATUS_CLOSED, string(got))
}
