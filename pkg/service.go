package mergeban

import (
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

type BanHammer interface {
	Ban(w http.ResponseWriter, r *http.Request)
	Lift(w http.ResponseWriter, r *http.Request)
}

func (b *banService) Ban(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		b.logger.Printf("Failed to parse form data from POST request: %v\n", err)
	}

	userID := r.FormValue("user_id")

	w.WriteHeader(200)
	_, err = w.Write([]byte(b.enqueueMerge(userID)))

	if err != nil {
		b.logger.Printf("Failed to write response body: %v\n", err)
	}
}

func (b *banService) Lift(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		b.logger.Printf("Failed to parse form data from POST request: %v\n", err)
	}

	userID := r.FormValue("user_id")

	w.WriteHeader(200)
	_, err = w.Write([]byte(b.withdrawMerge(userID)))

	if err != nil {
		b.logger.Printf("Failed to write response body: %v\n", err)
	}
}

func (b *banService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.Contains(path, "ban") {
		b.Ban(w, r)
	} else if strings.Contains(path, "lift") {
		b.Lift(w, r)
	} else {
		w.WriteHeader(404)

		_, err := w.Write([]byte("Unrecognized command"))
		if err != nil {
			b.logger.Printf("Failed to write 404 response body: %v\n", err)
		}
	}
}

func (b *banService) enqueueMerge(userID string) string {
	originalLength := b.mergeQueue.Length()

	b.mergeQueue.Enqueue(userID)

	if b.mergeQueue.Length() == 1 {
		return "You have banned merges!"
	} else if b.mergeQueue.Length() == originalLength {
		usersCurrentPosition := b.mergeQueue.FindIndex(userID) + 1
		return fmt.Sprintf("You are already in line! Your position: [%v/%v]\n", usersCurrentPosition, b.mergeQueue.Length())
	}

	return fmt.Sprintf("Hold up - someone else has banned merges. We'll message you when it's your turn to merge. Your position in line: [%v/%v]\n", b.mergeQueue.Length(), b.mergeQueue.Length())
}

func (b *banService) withdrawMerge(userID string) string {
	withdrawnUserID := b.mergeQueue.Withdraw(userID)

	if withdrawnUserID == nil {
		return "You aren't in line!"
	}

	return "You are no longer waiting to merge."
}
