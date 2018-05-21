package todo

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

//Testmain resets the table 'todos' in the test database before and after the suite of tests
func TestMain(m *testing.M) {
	db := OpenTestDB()
	db.Exec(DROP_TABLE)
	defer db.Exec(DROP_TABLE)
	db.Exec(RECREATE_SCHEMA)
	os.Exit(m.Run())
}

func uncheckedMarshal(i interface{}) []byte {
	bytes, _ := json.Marshal(i)
	return bytes
}

func byteReader(i interface{}) io.Reader {
	return bytes.NewReader(uncheckedMarshal(i))
}

func Test_create(t *testing.T) {
	wantOK := CreateOrUpdateTodo{"working!", STATUS_NEW}
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
			gotResp, err := create(OpenTestDB(), bytes.NewReader(tt.body))
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
	got, err := list(OpenTestDB())
	assert.NotNil(t, got)
	assert.NoError(t, err)
}

func Test_update(t *testing.T) {

	resp, _ := create(OpenTestDB(), byteReader(CreateOrUpdateTodo{"update working", STATUS_NEW}))
	var newTodo Todo
	json.Unmarshal(resp, &newTodo)

	want := Todo{newTodo.ID, "update working", STATUS_IN_PROGRESS}
	updateBody := CreateOrUpdateTodo{"update working", STATUS_IN_PROGRESS}

	gotResp, err := update(OpenTestDB(), byteReader(updateBody), strconv.FormatInt(int64(want.ID), 10))
	assert.NoError(t, err)

	var got Todo
	json.Unmarshal(gotResp, &got)
	assert.Equal(t, want, got)

	_, err = update(OpenTestDB(), byteReader(updateBody), "not parsable")
	assert.Error(t, err)

}
