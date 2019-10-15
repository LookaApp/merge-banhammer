package mergeban

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type banService struct {
	logger     *log.Logger
	mergeQueue []string
}

func CreateBanService() *banService {
	return &banService{
		logger:     log.New(os.Stdout, "", log.Ldate|log.Ltime),
		mergeQueue: make([]string, 0, 12),
	}
}

type BanHammer interface {
	Ban(w *http.ResponseWriter, r *http.Request)
}

func (b *banService) Ban(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	_, err := w.Write([]byte(b.EnqueueMerge()))

	if err != nil {
		b.logger.Printf("Failed to write response body: %v\n", err)
	}
}

func (b *banService) EnqueueMerge() string {
	b.mergeQueue = append(b.mergeQueue, "foobar")
	queueLength := len(b.mergeQueue)

	if queueLength == 1 {
		return "You have the merge banhammer!"
	}

	return fmt.Sprintf("The banhammer has already been taken. Your position: [%v/%v]\n", queueLength, queueLength)
}
