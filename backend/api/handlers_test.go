package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marcoswferreira/lofi-server/backend/models"
	"github.com/marcoswferreira/lofi-server/backend/service"
)

func TestGetStations(t *testing.T) {
	st := newMockStore()
	sm := service.NewStationManager()
	auth := service.NewAuthService(st)
	h := NewHandlers(sm, auth, st)

	sm.AddStation(models.Station{ID: "global", Name: "Global"})
	
	t.Run("Anonymous", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/stations", nil)
		rr := httptest.NewRecorder()
		h.GetStations(rr, req)
		var stations []models.Station
		json.NewDecoder(rr.Body).Decode(&stations)
		if len(stations) != 1 {
			t.Errorf("Expected 1 station, got %d", len(stations))
		}
	})

	t.Run("AuthenticatedWithShared", func(t *testing.T) {
		res, _ := auth.Register(context.Background(), models.RegisterRequest{
			Username: "user1", Email: "u1@test.com", Password: "pwd",
		})
		
		// Create a shared station in mock
		st.shares[1] = models.PlaylistShare{StationID: "shared1", UserID: res.User.ID, Status: "accepted"}
		sm.AddStation(models.Station{ID: "shared1", Name: "Shared"})

		req := httptest.NewRequest("GET", "/stations", nil)
		req.Header.Set("Authorization", "Bearer "+res.Token)
		rr := httptest.NewRecorder()
		h.GetStations(rr, req)

		var stations []models.Station
		json.NewDecoder(rr.Body).Decode(&stations)
		if len(stations) != 2 {
			t.Errorf("Expected 2 stations, got %d", len(stations))
		}
	})
}

func TestCreateStation(t *testing.T) {
	sm := service.NewStationManager()
	h := NewHandlers(sm, nil, nil)

	stationReq := models.CreateStationRequest{
		Name:        "New Station",
		Description: "Cool vibe",
		Playlist:    []models.Track{{ID: "t1", Title: "Track 1", Duration: time.Minute}},
	}
	body, _ := json.Marshal(stationReq)

	req := httptest.NewRequest("POST", "/stations", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.CreateStation(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetStationState(t *testing.T) {
	sm := service.NewStationManager()
	sm.AddStation(models.Station{
		ID:   "s1",
		Name: "S1",
		Playlist: []models.Track{{ID: "t1", Title: "T1", Duration: time.Minute}},
	})
	
	h := NewHandlers(sm, nil, nil)

	req := httptest.NewRequest("GET", "/stations/s1/state", nil)
	req.SetPathValue("id", "s1")
	rr := httptest.NewRecorder()

	h.GetStationState(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAuthHandlers(t *testing.T) {
	st := newMockStore()
	auth := service.NewAuthService(st)
	h := NewHandlers(service.NewStationManager(), auth, st)

	t.Run("Register", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.RegisterRequest{
			Username: "apiuser",
			Email:    "api@test.com",
			Password: "password",
		})
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()
		h.Register(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rr.Code)
		}
	})

	t.Run("Login", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.LoginRequest{
			Email:    "api@test.com",
			Password: "password",
		})
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()
		h.Login(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rr.Code)
		}
	})

	t.Run("LoginWrongPassword", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.LoginRequest{
			Email:    "api@test.com",
			Password: "wrong",
		})
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()
		h.Login(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rr.Code)
		}
	})

	t.Run("RegisterInvalidBody", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte("invalid")))
		rr := httptest.NewRecorder()
		h.Register(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", rr.Code)
		}
	})
}

func TestStationManagement(t *testing.T) {
	st := newMockStore()
	sm := service.NewStationManager()
	auth := service.NewAuthService(st)
	h := NewHandlers(sm, auth, st)

	res, _ := auth.Register(context.Background(), models.RegisterRequest{
		Username: "owner", Email: "owner@test.com", Password: "pwd",
	})
	token := res.Token

	t.Run("UpdateSuccess", func(t *testing.T) {
		sm.AddStation(models.Station{ID: "s1", Name: "S1", OwnerID: &res.User.ID})
		updateReq := models.CreateStationRequest{Name: "Updated"}
		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PUT", "/stations/s1", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.SetPathValue("id", "s1")
		rr := httptest.NewRecorder()
		h.UpdateStation(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Update failed: %d", rr.Code)
		}
	})

	t.Run("DeleteSuccess", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/stations/s1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.SetPathValue("id", "s1")
		rr := httptest.NewRecorder()
		h.DeleteStation(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Errorf("Delete failed: %d", rr.Code)
		}
	})

	t.Run("UpdateNotFound", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/stations/nonexistent", bytes.NewBuffer([]byte("{}")))
		req.SetPathValue("id", "nonexistent")
		rr := httptest.NewRecorder()
		h.UpdateStation(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", rr.Code)
		}
	})
}

func TestSharingHandlers(t *testing.T) {
	st := newMockStore()
	auth := service.NewAuthService(st)
	sm := service.NewStationManager()
	h := NewHandlers(sm, auth, st)

	resOwner, _ := auth.Register(context.Background(), models.RegisterRequest{
		Username: "owner", Email: "owner@test.com", Password: "hash",
	})
	resTarget, _ := auth.Register(context.Background(), models.RegisterRequest{
		Username: "target", Email: "target@test.com", Password: "hash",
	})
	
	token := resOwner.Token
	targetToken := resTarget.Token
	ownerID := resOwner.User.ID

	sm.AddStation(models.Station{ID: "s1", Name: "S1", OwnerID: &ownerID})

	t.Run("ShareSuccess", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.SharePlaylistRequest{UserEmail: "target@test.com"})
		req := httptest.NewRequest("POST", "/stations/s1/share", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.SetPathValue("id", "s1")
		rr := httptest.NewRecorder()
		h.SharePlaylist(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rr.Code)
		}
	})

	t.Run("GetInvitations", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/invitations", nil)
		req.Header.Set("Authorization", "Bearer "+targetToken)
		rr := httptest.NewRecorder()
		h.GetInvitations(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rr.Code)
		}
	})

	t.Run("UpdateInvitationSuccess", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.UpdateShareRequest{Status: "accepted"})
		req := httptest.NewRequest("PUT", "/auth/invitations/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+targetToken)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()
		h.UpdateInvitationStatus(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rr.Code)
		}
	})

	t.Run("SharePlaylistNotFound", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.SharePlaylistRequest{UserEmail: "nonexistent@test.com"})
		req := httptest.NewRequest("POST", "/stations/s1/share", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.SetPathValue("id", "s1")
		rr := httptest.NewRecorder()
		h.SharePlaylist(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", rr.Code)
		}
	})
}
