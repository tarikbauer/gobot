package app

import (
	"net/http"

	"github.com/tarikbauer/gobot/application"
)

type serverHandler struct {
	service application.Bot
}

func (sh *serverHandler) serve(w http.ResponseWriter, r *http.Request) {
	content, err := sh.service.Render(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	if content == nil {
		content = []byte("no data")
	}
	_, _ = w.Write(content)
}

func (sh *serverHandler) notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("endpoint not found"))
}

func (sh *serverHandler) notAllowed(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte("method not allowed"))
}
