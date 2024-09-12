package handler

import (
	"log/slog"
	"net/http"

	"github.com/diianpro/tingerDog/domain"
	"github.com/diianpro/tingerDog/service"
	"github.com/diianpro/tingerDog/transport/handler/utils"
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
		utils.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	utils.Render(w, r, &domain.ResponseUsers{Users: users})
}
