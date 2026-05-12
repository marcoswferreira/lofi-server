package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcoswferreira/lofi-server/backend/models"
)

type StoreInterface interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, string, error)
	SaveStation(ctx context.Context, station models.Station) error
	DeleteStation(ctx context.Context, id string) error
	CreateShare(ctx context.Context, stationID string, userID int) error
	UpdateShareStatus(ctx context.Context, shareID int, status string) error
	GetPendingShares(ctx context.Context, userID int) ([]models.PlaylistShare, error)
	GetAcceptedStationIDs(ctx context.Context, userID int) ([]string, error)
}

type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(ctx context.Context) (*Store, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/lofi_server"
	}

	var pool *pgxpool.Pool
	var err error

	// Retry connection a few times for Docker Compose startup
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				break
			}
		}
		log.Printf("Waiting for database... (attempt %d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database after retries: %w", err)
	}

	return &Store{Pool: pool}, nil
}

func (s *Store) InitSchema(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			is_premium BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS stations (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			playlist JSONB NOT NULL,
			owner_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS playlist_shares (
			id SERIAL PRIMARY KEY,
			station_id TEXT REFERENCES stations(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			status TEXT DEFAULT 'pending', -- pending, accepted, rejected
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(station_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS pomodoro_sessions (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			duration_minutes INTEGER NOT NULL,
			completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		if _, err := s.Pool.Exec(ctx, q); err != nil {
			return err
		}
	}
	return nil
}
