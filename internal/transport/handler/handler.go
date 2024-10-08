package handler

import (
	"log/slog"
	"net/http"

	"github.com/diianpro/tingerDog/internal/domain"
	"github.com/diianpro/tingerDog/internal/service"
	utils2 "github.com/diianpro/tingerDog/internal/transport/handler/utils"
)

type Handler struct {
	src *service.Service
}

func New(src *service.Service) *Handler {
	return &Handler{
		src: src,
	}
}

// GetAllUsers
// @Tags users
// @Router /list/user [get]
// @Success 200 {object} []domain.ResponseUsers "Get data successfully"
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.src.GetAllUsers(r.Context())
	if err != nil {
		slog.Error("failed to get all games")
		utils2.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	utils2.Render(w, r, &domain.ResponseUsers{Users: users})
}
