package server

import (
	"net/http"

	"github.com/extreme-business/lingo/apps/cms/views"
)

func NotFoundHandler(w http.ResponseWriter, vw views.Writer) {
	w.WriteHeader(http.StatusNotFound)
	vw.Error(w, "Page not found")
}
