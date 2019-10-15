package mergeban

import (
	"net/http"
)

type banService struct {
}

func CreateBanService() *banService {
	return &banService{}
}

type BanHammer interface {
	Ban(w *http.ResponseWriter, r *http.Request)
}

func (b *banService) Ban(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello World!"))
}
