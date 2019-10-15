package mergeban

import (
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
	_, err := w.Write([]byte("You have the merge banhammer!"))

	if err != nil {
		b.logger.Printf("Failed to write response body: %v\n", err)
	}
}
