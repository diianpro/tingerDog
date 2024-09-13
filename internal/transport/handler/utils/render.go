package utils

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func Render(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	w.Header().Set("Content-Type", "application/json")
	err := render.Render(w, r, v)
	if err != nil {
		slog.Error("failed render response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
