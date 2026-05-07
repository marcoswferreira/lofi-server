package store

import (
	"context"
	"encoding/json"

	"github.com/marcoswferreira/lofi-server/backend/models"
)

func (s *Store) SaveStation(ctx context.Context, station models.Station) error {
	playlistJSON, err := json.Marshal(station.Playlist)
	if err != nil {
		return err
	}

	_, err = s.Pool.Exec(ctx,
		`INSERT INTO stations (id, name, description, playlist, owner_id) 
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (id) DO UPDATE SET name = $2, description = $3, playlist = $4`,
		station.ID, station.Name, station.Description, playlistJSON, station.OwnerID,
	)
	return err
}

func (s *Store) DeleteStation(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx, "DELETE FROM stations WHERE id = $1", id)
	return err
}

func (s *Store) CreateShare(ctx context.Context, stationID string, userID int) error {
	_, err := s.Pool.Exec(ctx,
		"INSERT INTO playlist_shares (station_id, user_id, status) VALUES ($1, $2, 'pending') ON CONFLICT DO NOTHING",
		stationID, userID,
	)
	return err
}

func (s *Store) UpdateShareStatus(ctx context.Context, shareID int, status string) error {
	_, err := s.Pool.Exec(ctx, "UPDATE playlist_shares SET status = $1 WHERE id = $2", status, shareID)
	return err
}

func (s *Store) GetPendingShares(ctx context.Context, userID int) ([]models.PlaylistShare, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT ps.id, ps.station_id, s.name, ps.user_id, ps.status 
		 FROM playlist_shares ps 
		 JOIN stations s ON ps.station_id = s.id 
		 WHERE ps.user_id = $1 AND ps.status = 'pending'`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []models.PlaylistShare
	for rows.Next() {
		var ps models.PlaylistShare
		if err := rows.Scan(&ps.ID, &ps.StationID, &ps.StationName, &ps.UserID, &ps.Status); err != nil {
			return nil, err
		}
		shares = append(shares, ps)
	}
	return shares, nil
}

func (s *Store) GetAcceptedStationIDs(ctx context.Context, userID int) ([]string, error) {
	rows, err := s.Pool.Query(ctx, "SELECT station_id FROM playlist_shares WHERE user_id = $1 AND status = 'accepted'", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
