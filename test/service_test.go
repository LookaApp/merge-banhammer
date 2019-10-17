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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mergeban/pkg"
	"mergeban/test/mock_mergeban"
)

func TestBanEndpoint(t *testing.T) {
	t.Run("/ban - successfully acquiring initial lock", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		notifier := mock_mergeban.NewMockNotifier(mockCtrl)
		banService := mergeban.CreateService(logger, notifier)
		w, request := createBanRequest("42")

		banService.ServeHTTP(w, request)
		response := w.Result()
		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Contains(t, string(responseBody), "Stan Depain is waiting to merge!")
	})

	t.Run("/ban - queueing to acquire lock if it has already been taken", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		notifier := mock_mergeban.NewMockNotifier(mockCtrl)
		banService := mergeban.CreateService(logger, notifier)
		w, request := createBanRequest("23")
		w2, request2 := createBanRequest("42")

		banService.ServeHTTP(w, request)
		banService.ServeHTTP(w2, request2)
		response2 := w2.Result()
		responseBody2, _ := ioutil.ReadAll(response2.Body)

		assert.Equal(t, 200, response2.StatusCode)
		assert.Contains(t, string(responseBody2), "Hold up - someone else has banned merges. We'll message you when it's your turn to merge. Your position in line:")
		assert.Contains(t, string(responseBody2), "2/2")
	})

	t.Run("/ban - preventing the same user from enqueuing twice", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		notifier := mock_mergeban.NewMockNotifier(mockCtrl)
		banService := mergeban.CreateService(logger, notifier)
		w, request := createBanRequest("23")
		w2, request2 := createBanRequest("42")
		w3, request3 := createBanRequest("23")

		banService.ServeHTTP(w, request)
		banService.ServeHTTP(w2, request2)
		banService.ServeHTTP(w3, request3)
		response := w3.Result()

		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Contains(t, string(responseBody), "You are already in line! Your position:")
		assert.Contains(t, string(responseBody), "1/2")
	})
}

func TestLiftEndpoint(t *testing.T) {
	t.Run("/lift - warns when the lock is not held by this user", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		notifier := mock_mergeban.NewMockNotifier(mockCtrl)
		banService := mergeban.CreateService(logger, notifier)
		w, request := createLiftRequest("42")

		banService.ServeHTTP(w, request)
		response := w.Result()
		responseBody, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, 200, response.StatusCode)
		assert.Contains(t, string(responseBody), "You aren't in line!")
	})

	t.Run("/lift - releases any locks held by this user", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		notifier := mock_mergeban.NewMockNotifier(mockCtrl)
		notifier.
			EXPECT().
			Notify(gomock.Any(), gomock.Any()).
			AnyTimes()
		banService := mergeban.CreateService(logger, notifier)
		wBan, banRequest := createBanRequestWithResponseURL("42", "")
		wLift, liftRequest := createLiftRequest("42")

		banService.ServeHTTP(wBan, banRequest)
		banService.ServeHTTP(wLift, liftRequest)
		liftResponse := wLift.Result()
		responseBody, _ := ioutil.ReadAll(liftResponse.Body)

		assert.Equal(t, 200, liftResponse.StatusCode)
		assert.Contains(t, string(responseBody), "You are no longer waiting to merge.")
	})

	t.Run("/lift - notifies the next user in line if the user lifting held the lock", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
		notifier := mock_mergeban.NewMockNotifier(mockCtrl)
		banService := mergeban.CreateService(logger, notifier)
		wBan, banRequest := createBanRequestWithResponseURL("1", "http://example.com/respond/1")
		wBan2, banRequest2 := createBanRequestWithResponseURL("2", "http://example.com/respond/2")
		wLift, liftRequest := createLiftRequest("1")
		notifier.
			EXPECT().
			Notify(gomock.Eq("http://example.com/respond/2"), gomock.Eq("It's your turn to merge!")).
			Times(1)

		banService.ServeHTTP(wBan, banRequest)
		banService.ServeHTTP(wBan2, banRequest2)
		banService.ServeHTTP(wLift, liftRequest)
		liftResponse := wLift.Result()
		responseBody, _ := ioutil.ReadAll(liftResponse.Body)

		assert.Equal(t, 200, liftResponse.StatusCode)
		assert.Contains(t, string(responseBody), "You are no longer waiting to merge.")
	})

	/*
		t.Run("/bans - responds with a list of no users if the queue is empty", func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
			notifier := mock_mergeban.NewMockNotifier(mockCtrl)
			banService := mergeban.CreateService(logger, notifier)
			w, request := createListBansRequest()

			banService.ServeHTTP(w, request)
			response := w.Result()
			responseBody, _ := ioutil.ReadAll(response.Body)

			assert.Equal(t, 200, response.StatusCode)
			assert.Contains(t, string(responseBody), "Yay! There is no one currently waiting to merge.")
		})
	*/
}

func createLiftRequest(userID string) (*httptest.ResponseRecorder, *http.Request) {
	requestBody := strings.NewReader(fmt.Sprintf("user_id=%v", userID))
	request := httptest.NewRequest("POST", "/lift", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	return w, request
}

func createBanRequest(userID string) (*httptest.ResponseRecorder, *http.Request) {
	requestBody := strings.NewReader(fmt.Sprintf("user_id=%v&user_name=Stan Depain", userID))
	request := httptest.NewRequest("POST", "/ban", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	return w, request
}

func createBanRequestWithResponseURL(userID string, responseURL string) (*httptest.ResponseRecorder, *http.Request) {
	requestBody := strings.NewReader(fmt.Sprintf("user_id=%v&response_url=%v", userID, responseURL))
	request := httptest.NewRequest("POST", "/ban", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	return w, request
}

func createListBansRequest() (*httptest.ResponseRecorder, *http.Request) {
	requestBody := strings.NewReader("")
	request := httptest.NewRequest("POST", "/bans", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	return w, request
}
