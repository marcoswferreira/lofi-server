package api

import (
	"context"
	"errors"

	"github.com/marcoswferreira/lofi-server/backend/models"
)

type MockStore struct {
	users     map[string]*models.User
	passwords map[string]string
	stations  map[string]models.Station
	shares    map[int]models.PlaylistShare
	nextShareID int
}

func (m *MockStore) CreateUser(ctx context.Context, username, email, passwordHash string) (*models.User, error) {
	user := &models.User{ID: len(m.users) + 1, Username: username, Email: email}
	m.users[email] = user
	m.passwords[email] = passwordHash
	return user, nil
}

func (m *MockStore) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	u, ok := m.users[email]
	if !ok { return nil, "", errors.New("not found") }
	return u, m.passwords[email], nil
}

func (m *MockStore) SaveStation(ctx context.Context, station models.Station) error {
	m.stations[station.ID] = station
	return nil
}

func (m *MockStore) DeleteStation(ctx context.Context, id string) error {
	delete(m.stations, id)
	return nil
}

func (m *MockStore) CreateShare(ctx context.Context, stationID string, userID int) error {
	m.nextShareID++
	m.shares[m.nextShareID] = models.PlaylistShare{
		ID: m.nextShareID,
		StationID: stationID,
		UserID: userID,
		Status: "pending",
	}
	return nil
}

func (m *MockStore) UpdateShareStatus(ctx context.Context, shareID int, status string) error {
	s, ok := m.shares[shareID]
	if !ok { return errors.New("not found") }
	s.Status = status
	m.shares[shareID] = s
	return nil
}

func (m *MockStore) GetPendingShares(ctx context.Context, userID int) ([]models.PlaylistShare, error) {
	var res []models.PlaylistShare
	for _, s := range m.shares {
		if s.UserID == userID && s.Status == "pending" {
			res = append(res, s)
		}
	}
	return res, nil
}

func (m *MockStore) GetAcceptedStationIDs(ctx context.Context, userID int) ([]string, error) {
	var res []string
	for _, s := range m.shares {
		if s.UserID == userID && s.Status == "accepted" {
			res = append(res, s.StationID)
		}
	}
	return res, nil
}

func newMockStore() *MockStore {
	return &MockStore{
		users:     make(map[string]*models.User),
		passwords: make(map[string]string),
		stations:  make(map[string]models.Station),
		shares:    make(map[int]models.PlaylistShare),
	}
}
