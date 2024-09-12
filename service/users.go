package service

import (
	"context"

	"github.com/diianpro/tingerDog/domain"
)

func (s *Service) GetAllUsers(ctx context.Context) ([]domain.UserInfo, error) {
	var usersInfo []domain.UserInfo

	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		usersInfo = append(usersInfo, domain.UserInfo{
			Name: user.Name,
		})
	}

	return usersInfo, nil
}
