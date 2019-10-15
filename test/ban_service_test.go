package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"mergeban/mergeban"
)

func TestBanEndpoint(t *testing.T) {
	t.Run("/ban - successfully acquiring initial lock", func(t *testing.T) {
		banService := mergeban.CreateBanService()
		w, request := createBanRequest("42")

		banService.Ban(w, request)
		response := w.Result()
		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "You have the merge banhammer!", string(responseBody))
	})

	t.Run("/ban - queueing to acquire lock if it has already been taken", func(t *testing.T) {
		banService := mergeban.CreateBanService()
		w, request := createBanRequest("23")
		w2, request2 := createBanRequest("42")

		banService.Ban(w, request)
		banService.Ban(w2, request2)
		response2 := w2.Result()
		responseBody2, _ := ioutil.ReadAll(response2.Body)

		assert.Equal(t, 200, response2.StatusCode)
		assert.Equal(t, "The banhammer has already been taken. Your position: [2/2]\n", string(responseBody2))
	})

	t.Run("/ban - preventing the same user from enqueuing twice", func(t *testing.T) {
		banService := mergeban.CreateBanService()
		w, request := createBanRequest("23")
		w2, request2 := createBanRequest("42")
		w3, request3 := createBanRequest("23")

		banService.Ban(w, request)
		banService.Ban(w2, request2)
		banService.Ban(w3, request3)
		response := w3.Result()

		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "You are already in line! Your position: [1/2]\n", string(responseBody))
	})
}

func createBanRequest(userID string) (*httptest.ResponseRecorder, *http.Request) {
	requestBody := strings.NewReader(fmt.Sprintf("user_id=%v", userID))
	request := httptest.NewRequest("POST", "/ban", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	return w, request
}
