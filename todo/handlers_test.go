package todo

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const (
	DROP_TODOS      = "DROP TABLE TODO"
	RECREATE_SCHEMA = `CREATE TABLE public.todo (
	id serial NOT NULL,
	title varchar NULL,
	status varchar NULL,
	CONSTRAINT todo_pk PRIMARY KEY (id)
)
WITH (
	OIDS=FALSE
);

CREATE SEQUENCE IF NOT EXISTS public.todo_id_seq
NO MINVALUE
NO MAXVALUE;`
)

func TestMain(m *testing.M) {
	db := testDB()
	db.Exec(DROP_TODOS)
	defer db.Exec(DROP_TODOS)
	db.Exec(RECREATE_SCHEMA)
	os.Exit(m.Run())
}
func testDB() *sql.DB {
	db, _ := sql.Open("postgres",
		fmt.Sprintf("user=%s dbname=test sslmode=disable", os.Getenv("DB_USER")))
	return db
}
func uncheckedMarshal(i interface{}) []byte {
	bytes, _ := json.Marshal(i)
	return bytes
}

func byteReader(i interface{}) io.Reader {
	return bytes.NewReader(uncheckedMarshal(i))
}

func Test_create(t *testing.T) {
	wantOK := CreateOrUpdateTodo{"working!", "New"}
	badStatus := uncheckedMarshal(CreateOrUpdateTodo{"working!", "badStatus"})
	tests := []struct {
		name    string
		body    []byte
		want    CreateOrUpdateTodo
		wantErr bool
	}{
		{"ok", uncheckedMarshal(wantOK), wantOK, false},
		{"bad body", nil, CreateOrUpdateTodo{}, true},
		{"bad status", badStatus, CreateOrUpdateTodo{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := create(testDB(), bytes.NewReader(tt.body))
			if (err != nil) != tt.wantErr {
				t.Errorf("create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var got Todo
			json.Unmarshal(gotResp, &got)
			assert.Equal(t, tt.want.Status, got.Status)
			assert.Equal(t, tt.want.Title, got.Title)
		})
	}
}

func Test_list(t *testing.T) {
	got, err := list(testDB())
	assert.NotNil(t, got)
	assert.NoError(t, err)
}

func Test_update(t *testing.T) {

	resp, _ := create(testDB(), byteReader(CreateOrUpdateTodo{"update working", "New"}))
	var newTodo Todo
	json.Unmarshal(resp, &newTodo)

	want := Todo{newTodo.ID, "update working", STATUS_IN_PROGRESS}
	updateBody := CreateOrUpdateTodo{"update working", STATUS_IN_PROGRESS}

	gotResp, err := update(testDB(), byteReader(updateBody), strconv.FormatInt(int64(want.ID), 10))
	assert.NoError(t, err)

	var got Todo
	json.Unmarshal(gotResp, &got)
	assert.Equal(t, want, got)

	_, err = update(testDB(), byteReader(updateBody), "not parsable")
	assert.Error(t, err)
}
