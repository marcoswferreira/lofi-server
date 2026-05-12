package service

import (
	"sync"
	"time"

	"github.com/marcoswferreira/lofi-server/backend/models"
)

type StationManager struct {
	mu       sync.RWMutex
	stations map[string]*stationRuntime
}

type stationRuntime struct {
	station      models.Station
	currentIndex int
	startTime    time.Time
}

func NewStationManager() *StationManager {
	return &StationManager{
		stations: make(map[string]*stationRuntime),
	}
}

func (sm *StationManager) AddStation(station models.Station) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.stations[station.ID] = &stationRuntime{
		station:      station,
		currentIndex: 0,
		startTime:    time.Now(),
	}
}

func (sm *StationManager) GetState(stationID string) (models.StationState, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	runtime, ok := sm.stations[stationID]
	if !ok {
		return models.StationState{}, false
	}

	// Update current track if it has finished
	sm.updateStation(runtime)

	currentTrack := runtime.station.Playlist[runtime.currentIndex]
	now := time.Now()
	elapsed := now.Sub(runtime.startTime)

	return models.StationState{
		StationID:      stationID,
		CurrentTrack:   currentTrack,
		StartTime:      runtime.startTime.Unix(),
		CurrentSeconds: int64(elapsed.Seconds()),
	}, true
}

func (sm *StationManager) updateStation(runtime *stationRuntime) {
	if len(runtime.station.Playlist) == 0 {
		return
	}

	now := time.Now()
	for {
		currentTrack := runtime.station.Playlist[runtime.currentIndex]
		if currentTrack.Duration <= 0 {
			// Prevent infinite loop: if duration is invalid, just stay on this track
			// but advance to next to try to find a valid one?
			// For now, let's just break to avoid hang.
			break
		}
		if now.Sub(runtime.startTime) < currentTrack.Duration {
			break
		}

		// Advance to next track
		runtime.startTime = runtime.startTime.Add(currentTrack.Duration)
		runtime.currentIndex = (runtime.currentIndex + 1) % len(runtime.station.Playlist)
	}
}

func (sm *StationManager) GetStationByID(id string) (models.Station, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	runtime, ok := sm.stations[id]
	if !ok {
		return models.Station{}, false
	}
	return runtime.station, true
}

func (sm *StationManager) GetStationsForUser(userID *int, sharedIDs []string) []models.Station {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Convert sharedIDs to a map for O(1) lookup
	sharedMap := make(map[string]bool)
	for _, id := range sharedIDs {
		sharedMap[id] = true
	}

	var stations []models.Station
	for _, runtime := range sm.stations {
		// Include global stations (OwnerID == nil)
		// OR private stations if the user is the owner
		// OR private stations if they are shared and accepted
		isOwner := userID != nil && runtime.station.OwnerID != nil && *runtime.station.OwnerID == *userID
		isShared := sharedMap[runtime.station.ID]

		if runtime.station.OwnerID == nil || isOwner || isShared {
			stations = append(stations, runtime.station)
		}
	}
	return stations
}

func (sm *StationManager) UpdateStation(id string, station models.Station) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	runtime, ok := sm.stations[id]
	if !ok {
		return false
	}

	runtime.station = station
	// We keep the current playback progress (currentIndex and startTime)
	// unless the playlist has changed significantly.
	// For simplicity, we just keep them.
	return true
}

func (sm *StationManager) DeleteStation(id string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, ok := sm.stations[id]; !ok {
		return false
	}

	delete(sm.stations, id)
	return true
}
