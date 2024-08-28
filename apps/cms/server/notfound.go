package server

import (
	"net/http"

	"github.com/extreme-business/lingo/apps/cms/views"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	views.Error(w, "Page not found")
}
