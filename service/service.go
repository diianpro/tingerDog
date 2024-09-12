package service

import "github.com/diianpro/tingerDog/repository"

type Service struct {
	userRepo repository.Repo
}

func New(userRepo repository.Repo) *Service {
	return &Service{userRepo}
}
