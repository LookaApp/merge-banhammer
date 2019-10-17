package mergeban

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type banService struct {
	logger     *log.Logger
	mergeQueue *mergeQueue
	notifier   Notifier
}

func CreateService(logger *log.Logger, notifier Notifier) *banService {
	return &banService{
		logger:     logger,
		mergeQueue: NewQueue(),
		notifier:   notifier,
	}
}

func marshalResponse(responseText string) ([]byte, error) {
	responsePayload := map[string]string{
		"response_type": "in_channel",
		"text":          responseText,
	}
	responseJSON, err := json.Marshal(responsePayload)
	if err != nil {
		return []byte{}, err
	}

	return responseJSON, nil
}

func (b *banService) Ban(responseURL, userID, userName string) ([]byte, error) {
	responseText := b.enqueueMerge(userID, responseURL, userName)

	return marshalResponse(responseText)
}

func (b *banService) Lift(userID string) ([]byte, error) {
	responseText := b.withdrawMerge(userID)

	return marshalResponse(responseText)
}

func (b *banService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		b.logger.Printf("Failed to parse form data from POST request: %v\n", err)
	}

	responseURL := r.FormValue("response_url")
	userID := r.FormValue("user_id")
	userName := r.FormValue("user_name")

	path := r.URL.Path
	if strings.Contains(path, "ban") {
		w.Header().Add("Content-Type", "application/json")

		responsePayload, err := b.Ban(responseURL, userID, userName)
		if err != nil {
			b.logger.Printf("Failed to handle ban request: %v\n", err)
			w.WriteHeader(500)
		}

		_, err = w.Write(responsePayload)
		if err != nil {
			b.logger.Printf("Failed to write response body: %v\n", err)
			w.WriteHeader(500)
		}

		w.WriteHeader(200)
	} else if strings.Contains(path, "lift") {
		w.Header().Add("Content-Type", "application/json")

		responsePayload, err := b.Lift(userID)
		if err != nil {
			b.logger.Printf("Failed to handle lift request: %v\n", err)
			w.WriteHeader(500)
		}

		_, err = w.Write(responsePayload)
		if err != nil {
			b.logger.Printf("Failed to write response body: %v\n", err)
			w.WriteHeader(500)
		}

		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)

		_, err := w.Write([]byte("Unrecognized command"))
		if err != nil {
			b.logger.Printf("Failed to write 404 response body: %v\n", err)
		}
	}
}

func (b *banService) enqueueMerge(userID, responseURL, userName string) string {
	originalLength := b.mergeQueue.Length()

	b.mergeQueue.Enqueue(userID, userName, responseURL)

	if b.mergeQueue.Length() == 1 {
		return fmt.Sprintf("%s is waiting to merge!", userName)
	} else if b.mergeQueue.Length() == originalLength {
		usersCurrentPosition := b.mergeQueue.FindIndex(userID) + 1
		return fmt.Sprintf("You are already in line! Your position: [%v/%v]\n", usersCurrentPosition, b.mergeQueue.Length())
	}

	return fmt.Sprintf("Hold up - someone else has banned merges. We'll message you when it's your turn to merge. Your position in line: [%v/%v]\n", b.mergeQueue.Length(), b.mergeQueue.Length())
}

func (b *banService) withdrawMerge(userID string) string {
	positionOfWithdrawingUser := b.mergeQueue.FindIndex(userID)

	if positionOfWithdrawingUser == -1 {
		return "You aren't in line!"
	} else if positionOfWithdrawingUser == 0 {
		b.mergeQueue.Dequeue()
		nextInLine := b.mergeQueue.Peek()
		if nextInLine != nil {
			b.notifier.Notify(nextInLine.ResponseURL, "It's your turn to merge!")
		}
	} else {
		b.mergeQueue.Withdraw(userID)
	}

	return "You are no longer waiting to merge."
}
