package test

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"mergeban/mergeban"
)

func TestBanEndpoint(t *testing.T) {
	t.Run("/ban - successfully acquiring initial lock", func(t *testing.T) {
		banService := mergeban.CreateBanService()
		requestBody := strings.NewReader("user_id=42")
		request := httptest.NewRequest("POST", "/ban", requestBody)
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		banService.Ban(w, request)
		response := w.Result()
		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "You have the merge banhammer!", string(responseBody))
	})

	t.Run("/ban - queueing to acquire lock if it has already been taken", func(t *testing.T) {
		banService := mergeban.CreateBanService()
		requestBody := strings.NewReader("user_id=23")
		request := httptest.NewRequest("POST", "/ban", requestBody)
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		requestBody2 := strings.NewReader("user_id=42")
		request2 := httptest.NewRequest("POST", "/ban", requestBody2)
		request2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w2 := httptest.NewRecorder()

		banService.Ban(w, request)
		banService.Ban(w2, request2)
		response2 := w2.Result()

		responseBody2, _ := ioutil.ReadAll(response2.Body)

		assert.Equal(t, 200, response2.StatusCode)
		assert.Equal(t, "The banhammer has already been taken. Your position: [2/2]\n", string(responseBody2))
	})

	t.Run("/ban - preventing the same user from enqueuing twice", func(t *testing.T) {
		banService := mergeban.CreateBanService()
		requestBody := strings.NewReader("user_id=23")
		request := httptest.NewRequest("POST", "/ban", requestBody)
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		requestBody2 := strings.NewReader("user_id=42")
		request2 := httptest.NewRequest("POST", "/ban", requestBody2)
		request2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w2 := httptest.NewRecorder()
		w3 := httptest.NewRecorder()

		banService.Ban(w, request)
		banService.Ban(w2, request2)
		banService.Ban(w3, request)
		response := w3.Result()

		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "You are already in line! Your position: [1/2]\n", string(responseBody))
	})
}
