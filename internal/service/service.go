package service

import (
	"github.com/diianpro/tingerDog/internal/storage/postgres"
)

type Service struct {
	userRepo postgres.Repo
}

func New(userRepo postgres.Repo) *Service {
	return &Service{userRepo}
}
