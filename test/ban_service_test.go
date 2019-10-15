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
	banService := mergeban.CreateBanService()

	t.Run("/ban", func(t *testing.T) {
		requestBody := strings.NewReader("user_id=42")
		request := httptest.NewRequest("POST", "/ban", requestBody)
		w := httptest.NewRecorder()

		banService.Ban(w, request)
		response := w.Result()
		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.NotEmpty(t, responseBody)
	})
}
