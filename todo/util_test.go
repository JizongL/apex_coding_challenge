package todo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_writeErr(t *testing.T) {
	rec := httptest.NewRecorder()
	writeErr(rec, BadRequestError("some bad request"))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "some bad request\n", rec.Body.String())

	rec = httptest.NewRecorder()
	internalErr := fmt.Errorf("some error")
	writeErr(rec, internalErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error\n", rec.Body.String())

}
