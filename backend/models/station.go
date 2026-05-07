package models

import "time"

type Track struct {
	ID       string        `json:"id"`       // YouTube Video ID
	Title    string        `json:"title"`
	Duration time.Duration `json:"duration" swaggertype:"primitive,integer"`
}

type Station struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Playlist    []Track `json:"playlist"`
	OwnerID     *int    `json:"ownerId,omitempty"` // nil means global/system station
}

type CreateStationRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Playlist    []Track `json:"playlist"`
}

type StationState struct {
	StationID      string `json:"stationId"`
	CurrentTrack   Track  `json:"currentTrack"`
	StartTime      int64  `json:"startTime"`      // Unix timestamp when the current track started
	CurrentSeconds int64  `json:"currentSeconds"` // Offset in seconds from start
}

type PlaylistShare struct {
	ID          int    `json:"id"`
	StationID   string `json:"stationId"`
	StationName string `json:"stationName,omitempty"`
	UserID      int    `json:"userId"`
	Status      string `json:"status"` // pending, accepted, rejected
}

type SharePlaylistRequest struct {
	UserEmail string `json:"userEmail"`
}

type UpdateShareRequest struct {
	Status string `json:"status"` // accepted, rejected
}
