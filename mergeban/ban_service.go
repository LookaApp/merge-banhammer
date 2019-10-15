package mergeban

import (
	"fmt"
	"log"
	"net/http"
)

type banService struct {
	logger     *log.Logger
	mergeQueue []string
}

func CreateBanService(logger *log.Logger) *banService {
	return &banService{
		logger:     logger,
		mergeQueue: make([]string, 0, 12),
	}
}

type BanHammer interface {
	Ban(w http.ResponseWriter, r *http.Request)
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

func (b *banService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.Ban(w, r)
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
