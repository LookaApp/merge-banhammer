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

	"mergeban/pkg"
)

func TestBanEndpoint(t *testing.T) {
	t.Run("/ban - successfully acquiring initial lock", func(t *testing.T) {
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		banService := mergeban.CreateService(logger)
		w, request := createBanRequest("42")

		banService.ServeHTTP(w, request)
		response := w.Result()
		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "You have the merge banhammer!", string(responseBody))
	})

	t.Run("/ban - queueing to acquire lock if it has already been taken", func(t *testing.T) {
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		banService := mergeban.CreateService(logger)
		w, request := createBanRequest("23")
		w2, request2 := createBanRequest("42")

		banService.ServeHTTP(w, request)
		banService.ServeHTTP(w2, request2)
		response2 := w2.Result()
		responseBody2, _ := ioutil.ReadAll(response2.Body)

		assert.Equal(t, 200, response2.StatusCode)
		assert.Equal(t, "The banhammer has already been taken. Your position: [2/2]\n", string(responseBody2))
	})

	t.Run("/ban - preventing the same user from enqueuing twice", func(t *testing.T) {
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		banService := mergeban.CreateService(logger)
		w, request := createBanRequest("23")
		w2, request2 := createBanRequest("42")
		w3, request3 := createBanRequest("23")

		banService.ServeHTTP(w, request)
		banService.ServeHTTP(w2, request2)
		banService.ServeHTTP(w3, request3)
		response := w3.Result()

		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "You are already in line! Your position: [1/2]\n", string(responseBody))
	})
}

func TestLiftEndpoint(t *testing.T) {
	t.Run("/lift - warns when the lock is not held by this user", func(t *testing.T) {
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		banService := mergeban.CreateService(logger)
		w, request := createLiftRequest("42")

		banService.ServeHTTP(w, request)
		response := w.Result()
		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "You do not have the banhammer!", string(responseBody))
	})

	t.Run("/lift - releases any locks held by this user", func(t *testing.T) {
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		banService := mergeban.CreateService(logger)

		wBan, banRequest := createBanRequest("42")
		wLift, liftRequest := createLiftRequest("42")

		banService.ServeHTTP(wBan, banRequest)
		banService.ServeHTTP(wLift, liftRequest)
		liftResponse := wLift.Result()
		responseBody, _ := ioutil.ReadAll(liftResponse.Body)

		assert.Equal(t, 200, liftResponse.StatusCode)
		assert.Equal(t, "You no longer have the banhammer!", string(responseBody))
	})
}

func createLiftRequest(userID string) (*httptest.ResponseRecorder, *http.Request) {
	requestBody := strings.NewReader(fmt.Sprintf("user_id=%v", userID))
	request := httptest.NewRequest("POST", "/lift", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	return w, request
}

func createBanRequest(userID string) (*httptest.ResponseRecorder, *http.Request) {
	requestBody := strings.NewReader(fmt.Sprintf("user_id=%v", userID))
	request := httptest.NewRequest("POST", "/ban", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	return w, request
}
