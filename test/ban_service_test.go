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
		w := httptest.NewRecorder()
		requestBody2 := strings.NewReader("user_id=42")
		request2 := httptest.NewRequest("POST", "/ban", requestBody2)
		w2 := httptest.NewRecorder()

		banService.Ban(w, request)
		banService.Ban(w2, request2)
		response2 := w2.Result()

		responseBody2, _ := ioutil.ReadAll(response2.Body)

		assert.Equal(t, 200, response2.StatusCode)
		assert.Equal(t, "The banhammer has already been taken. Your position: [2/2]\n", string(responseBody2))
	})
}
