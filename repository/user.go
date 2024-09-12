package repository

import (
	"context"

	"github.com/diianpro/tingerDog/repository/models"
)

func (r *Repository) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users`
	rows, err := r.db(ctx).Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
