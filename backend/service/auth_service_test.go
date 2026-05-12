package service

import (
	"context"
	"errors"
	"testing"

	"github.com/marcoswferreira/lofi-server/backend/models"
)

// MockStore implements store.StoreInterface for testing
type MockStore struct {
	users     map[string]*models.User // email -> user
	passwords map[string]string       // email -> passwordHash
	stations  map[string]models.Station
	shares    []models.PlaylistShare
	failNext  bool
}

func (m *MockStore) CreateUser(ctx context.Context, username, email, passwordHash string) (*models.User, error) {
	if m.failNext {
		return nil, errors.New("db error")
	}
	user := &models.User{
		ID:       len(m.users) + 1,
		Username: username,
		Email:    email,
	}
	m.users[email] = user
	m.passwords[email] = passwordHash
	return user, nil
}

func (m *MockStore) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, "", errors.New("not found")
	}
	return user, m.passwords[email], nil
}

func (m *MockStore) SaveStation(ctx context.Context, station models.Station) error       { return nil }
func (m *MockStore) DeleteStation(ctx context.Context, id string) error                  { return nil }
func (m *MockStore) CreateShare(ctx context.Context, stationID string, userID int) error { return nil }
func (m *MockStore) UpdateShareStatus(ctx context.Context, shareID int, status string) error {
	return nil
}
func (m *MockStore) GetPendingShares(ctx context.Context, userID int) ([]models.PlaylistShare, error) {
	return nil, nil
}
func (m *MockStore) GetAcceptedStationIDs(ctx context.Context, userID int) ([]string, error) {
	return nil, nil
}

func TestAuthService_ValidateToken(t *testing.T) {
	mock := &MockStore{}
	auth := NewAuthService(mock)

	token, _ := auth.generateToken(&models.User{ID: 123})

	t.Run("ValidToken", func(t *testing.T) {
		id, err := auth.ValidateToken(token)
		if err != nil || id != 123 {
			t.Errorf("Expected id 123, got %d, err: %v", id, err)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, err := auth.ValidateToken("invalid")
		if err == nil {
			t.Error("Expected error for invalid token")
		}
	})

	t.Run("WrongSecret", func(t *testing.T) {
		wrongAuth := NewAuthService(mock)
		wrongAuth.secret = []byte("wrong")
		_, err := wrongAuth.ValidateToken(token)
		if err == nil {
			t.Error("Expected error for token with wrong secret")
		}
	})
}

func TestAuthService_RegisterAndLogin(t *testing.T) {
	mock := &MockStore{
		users:     make(map[string]*models.User),
		passwords: make(map[string]string),
	}
	auth := NewAuthService(mock)
	ctx := context.Background()

	t.Run("RegisterSuccess", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		res, err := auth.Register(ctx, req)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
		if res.User.Username != "testuser" || res.Token == "" {
			t.Errorf("Invalid response: %v", res)
		}
	})

	t.Run("LoginSuccess", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		res, err := auth.Login(ctx, req)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if res.User.Email != "test@example.com" || res.Token == "" {
			t.Errorf("Invalid response: %v", res)
		}
	})

	t.Run("LoginWrongPassword", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrong",
		}
		_, err := auth.Login(ctx, req)
		if err == nil {
			t.Error("Expected error for wrong password")
		}
	})

	t.Run("LoginNotFound", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password",
		}
		_, err := auth.Login(ctx, req)
		if err == nil {
			t.Error("Expected error for non-existent user")
		}
	})
}
