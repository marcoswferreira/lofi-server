package store

import (
	"context"

	"github.com/marcoswferreira/lofi-server/backend/models"
)

func (s *Store) CreateUser(ctx context.Context, username, email, passwordHash string) (*models.User, error) {
	var user models.User
	err := s.Pool.QueryRow(ctx, 
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, username, email, is_premium, created_at",
		username, email, passwordHash,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsPremium, &user.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	var user models.User
	var passwordHash string
	err := s.Pool.QueryRow(ctx,
		"SELECT id, username, email, password_hash, is_premium, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.IsPremium, &user.CreatedAt)
	
	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}
