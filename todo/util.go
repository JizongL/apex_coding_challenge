package todo

import (
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
	fmt.Fprint(w, string(jsonResp))
}

type key string

const (
	ENV_USER key = "DB_USER"
	ENV_NAME key = "DB_NAME"
	ENV_HOST key = "DB_HOST"
	ENV_PASS key = "DB_PASSWORD"

	ENV_TEST_USER key = "TEST_DB_USER"
	ENV_TEST_NAME key = "TEST_DB_NAME"
)

func (k key) val() (string, error) {
	v, ok := os.LookupEnv(string(k))
	if !ok {
		return "", fmt.Errorf("missing environment key $%s", k)
	}
	return v, nil
}
