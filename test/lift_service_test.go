package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"mergeban/mergeban"
)

func TestLiftEndpoint(t *testing.T) {
	t.Run("/lift - successfully releasing a held lock", func(t *testing.T) {
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		banService := mergeban.CreateService(logger)
		w, request := createLiftRequest("42")

		banService.Lift(w, request)
		response := w.Result()
		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.NotEmpty(t, responseBody)
	})
}

func createLiftRequest(userID string) (*httptest.ResponseRecorder, *http.Request) {
	requestBody := strings.NewReader(fmt.Sprintf("user_id=%v", userID))
	request := httptest.NewRequest("POST", "/lift", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	return w, request
}
