package mergeban

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type banService struct {
	logger     *log.Logger
	mergeQueue []string
}

func CreateService(logger *log.Logger) *banService {
	return &banService{
		logger:     logger,
		mergeQueue: make([]string, 0, 12),
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

	_, err = w.Write([]byte(b.dequeueMerge(userID)))
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
	queueLength := len(b.mergeQueue)

	for position, enqueuedID := range b.mergeQueue {
		if enqueuedID == userID {
			return fmt.Sprintf("You are already in line! Your position: [%v/%v]\n", position+1, queueLength)
		}
	}

	b.mergeQueue = append(b.mergeQueue, userID)

	if queueLength == 0 {
		return "You have the merge banhammer!"
	}

	return fmt.Sprintf("The banhammer has already been taken. Your position: [%v/%v]\n", queueLength+1, queueLength+1)
}

func (b *banService) dequeueMerge(userID string) string {
	userPosition := int(-1)

	for position, enqueuedID := range b.mergeQueue {
		if enqueuedID == userID {
			userPosition = position
		}
	}

	if userPosition == -1 {
		return "You do not have the banhammer!"
	}

	copy(b.mergeQueue[userPosition:], b.mergeQueue[userPosition+1:])
	return "You no longer have the banhammer!"
}
