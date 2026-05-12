package service

import (
	"testing"
	"time"

	"github.com/marcoswferreira/lofi-server/backend/models"
)

func TestNewStationManager(t *testing.T) {
	sm := NewStationManager()
	if sm.stations == nil {
		t.Error("Expected stations map to be initialized")
	}
}

func TestStationLifecycle(t *testing.T) {
	sm := NewStationManager()
	station := models.Station{
		ID:   "test-id",
		Name: "Test Station",
		Playlist: []models.Track{
			{ID: "track-1", Title: "Track 1", Duration: time.Minute},
		},
	}

	// Add
	sm.AddStation(station)
	if s, ok := sm.GetStationByID("test-id"); !ok || s.Name != "Test Station" {
		t.Errorf("Expected station to be added and retrievable, got ok=%v, name=%s", ok, s.Name)
	}

	// Update
	station.Name = "Updated Station"
	sm.UpdateStation("test-id", station)
	if s, _ := sm.GetStationByID("test-id"); s.Name != "Updated Station" {
		t.Errorf("Expected station name to be updated, got %s", s.Name)
	}

	// Non-existent update
	if ok := sm.UpdateStation("invalid-id", station); ok {
		t.Error("Expected UpdateStation to return false for non-existent station")
	}

	// Delete
	if ok := sm.DeleteStation("test-id"); !ok {
		t.Error("Expected station to be deleted")
	}
	if _, ok := sm.GetStationByID("test-id"); ok {
		t.Error("Expected station to be gone after deletion")
	}

	// Non-existent delete
	if ok := sm.DeleteStation("invalid-id"); ok {
		t.Error("Expected DeleteStation to return false for non-existent station")
	}
}

func TestGetStationsForUser(t *testing.T) {
	sm := NewStationManager()
	owner1 := 1
	owner2 := 2

	sm.AddStation(models.Station{ID: "global", Name: "Global", OwnerID: nil})
	sm.AddStation(models.Station{ID: "owner1", Name: "Owner 1", OwnerID: &owner1})
	sm.AddStation(models.Station{ID: "owner2", Name: "Owner 2", OwnerID: &owner2})

	t.Run("AnonymousUser", func(t *testing.T) {
		stations := sm.GetStationsForUser(nil, nil)
		if len(stations) != 1 || stations[0].ID != "global" {
			t.Errorf("Anonymous user should only see global station, got %v", stations)
		}
	})

	t.Run("Owner1Only", func(t *testing.T) {
		stations := sm.GetStationsForUser(&owner1, nil)
		if len(stations) != 2 {
			t.Errorf("Owner1 should see global and their own station, got %d", len(stations))
		}
	})

	t.Run("Owner1WithSharedStation", func(t *testing.T) {
		stations := sm.GetStationsForUser(&owner1, []string{"owner2"})
		if len(stations) != 3 {
			t.Errorf("Owner1 with share should see global, their own, and shared station, got %d", len(stations))
		}
	})
}

func TestGetState(t *testing.T) {
	sm := NewStationManager()
	station := models.Station{
		ID:   "s1",
		Name: "Station 1",
		Playlist: []models.Track{
			{ID: "t1", Title: "Track 1", Duration: 2 * time.Second},
			{ID: "t2", Title: "Track 2", Duration: 2 * time.Second},
		},
	}
	sm.AddStation(station)

	t.Run("InitialState", func(t *testing.T) {
		state, ok := sm.GetState("s1")
		if !ok || state.CurrentTrack.ID != "t1" {
			t.Errorf("Expected Track 1 at start, got %v", state)
		}
	})

	t.Run("NonExistentStation", func(t *testing.T) {
		_, ok := sm.GetState("invalid")
		if ok {
			t.Error("Expected ok=false for non-existent station")
		}
	})

	t.Run("AdvancedState", func(t *testing.T) {
		// Mock time by adjusting the start time back manually
		sm.mu.Lock()
		sm.stations["s1"].startTime = time.Now().Add(-3 * time.Second)
		sm.stations["s1"].currentIndex = 0 // Reset index for consistent test
		sm.mu.Unlock()

		state, ok := sm.GetState("s1")
		if !ok || state.CurrentTrack.ID != "t2" {
			t.Errorf("Expected Track 2 after 3 seconds, got %v", state)
		}
		// Track 1 (2s) + 1s into Track 2 = 3s total, so CurrentSeconds (offset in current track) should be 1
		if state.CurrentSeconds != 1 {
			t.Errorf("Expected 1 elapsed second in current track, got %d", state.CurrentSeconds)
		}
	})

	t.Run("LoopAroundState", func(t *testing.T) {
		sm.mu.Lock()
		sm.stations["s1"].startTime = time.Now().Add(-5 * time.Second)
		sm.stations["s1"].currentIndex = 0 // Reset index for consistent test
		sm.mu.Unlock()

		state, ok := sm.GetState("s1")
		if !ok || state.CurrentTrack.ID != "t1" {
			t.Errorf("Expected Track 1 after 5 seconds (loop around), got %v", state)
		}
		if state.CurrentSeconds != 1 {
			t.Errorf("Expected 1 elapsed second after loop, got %d", state.CurrentSeconds)
		}
	})
}
