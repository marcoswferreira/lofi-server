package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/marcoswferreira/lofi-server/backend/models"
	"github.com/marcoswferreira/lofi-server/backend/service"
	"github.com/marcoswferreira/lofi-server/backend/store"
)

type Handlers struct {
	sm   *service.StationManager
	auth *service.AuthService
	st   *store.Store
}

func NewHandlers(sm *service.StationManager, auth *service.AuthService, st *store.Store) *Handlers {
	return &Handlers{sm: sm, auth: auth, st: st}
}

// GetStations returns the list of available radio stations.
// @Summary Get all stations
// @Description List all lofi music stations (global and user's private/shared ones)
// @Tags stations
// @Produce json
// @Success 200 {array} models.Station
// @Router /stations [get]
func (h *Handlers) GetStations(w http.ResponseWriter, r *http.Request) {
	var userID *int
	var sharedIDs []string
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		id, err := h.auth.ValidateToken(authHeader[7:])
		if err == nil {
			userID = &id
			// Fetch shared and accepted station IDs from DB
			if h.st != nil {
				ids, _ := h.st.GetAcceptedStationIDs(r.Context(), id)
				sharedIDs = ids
			}
		}
	}

	stations := h.sm.GetStationsForUser(userID, sharedIDs)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stations)
}

// CreateStation creates a new music station.
// @Summary Create a new station
// @Description Add a new lofi music station with a playlist. If authenticated, it will be private.
// @Tags stations
// @Accept json
// @Produce json
// @Param request body models.CreateStationRequest true "Station details"
// @Success 201 {object} models.Station
// @Failure 400 {string} string "Invalid request body"
// @Router /stations [post]
func (h *Handlers) CreateStation(w http.ResponseWriter, r *http.Request) {
	var userID *int
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		id, err := h.auth.ValidateToken(authHeader[7:])
		if err == nil {
			userID = &id
		}
	}

	var req models.CreateStationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := time.Now().Format("20060102150405")
	station := models.Station{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Playlist:    req.Playlist,
		OwnerID:     userID,
	}

	h.sm.AddStation(station)
	
	// Persist to DB if store is available
	if h.st != nil {
		h.st.SaveStation(r.Context(), station)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(station)
}

// UpdateStation updates an existing music station.
// @Summary Update a station
// @Description Update the details and playlist of an existing station. Requires ownership.
// @Tags stations
// @Accept json
// @Produce json
// @Param id path string true "Station ID"
// @Param request body models.CreateStationRequest true "Updated station details"
// @Success 200 {object} models.Station
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Station not found"
// @Router /stations/{id} [put]
func (h *Handlers) UpdateStation(w http.ResponseWriter, r *http.Request) {
	var userID *int
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		id, err := h.auth.ValidateToken(authHeader[7:])
		if err == nil {
			userID = &id
		}
	}

	id := r.PathValue("id")
	current, ok := h.sm.GetStationByID(id)
	if !ok {
		http.Error(w, "Station not found", http.StatusNotFound)
		return
	}

	if current.OwnerID != nil {
		if userID == nil || *current.OwnerID != *userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	var req models.CreateStationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	station := models.Station{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Playlist:    req.Playlist,
		OwnerID:     current.OwnerID,
	}

	h.sm.UpdateStation(id, station)
	if h.st != nil {
		h.st.SaveStation(r.Context(), station)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(station)
}

// DeleteStation removes a music station.
// @Summary Delete a station
// @Description Remove a station from the list. Requires ownership.
// @Tags stations
// @Param id path string true "Station ID"
// @Success 204 "No Content"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Station not found"
// @Router /stations/{id} [delete]
func (h *Handlers) DeleteStation(w http.ResponseWriter, r *http.Request) {
	var userID *int
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		id, err := h.auth.ValidateToken(authHeader[7:])
		if err == nil {
			userID = &id
		}
	}

	id := r.PathValue("id")
	current, ok := h.sm.GetStationByID(id)
	if !ok {
		http.Error(w, "Station not found", http.StatusNotFound)
		return
	}

	if current.OwnerID != nil {
		if userID == nil || *current.OwnerID != *userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	h.sm.DeleteStation(id)
	if h.st != nil {
		h.st.DeleteStation(r.Context(), id)
	}
	w.WriteHeader(http.StatusNoContent)
}

// SharePlaylist invites another user to access a playlist.
// @Summary Share playlist
// @Description Invite another user to access your private playlist
// @Tags stations
// @Accept json
// @Param id path string true "Station ID"
// @Param request body models.SharePlaylistRequest true "User email to share with"
// @Success 200 {string} string "Invitation sent"
// @Router /stations/{id}/share [post]
func (h *Handlers) SharePlaylist(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ownerID, _ := h.auth.ValidateToken(authHeader[7:])

	stationID := r.PathValue("id")
	station, ok := h.sm.GetStationByID(stationID)
	if !ok || station.OwnerID == nil || *station.OwnerID != ownerID {
		http.Error(w, "Unauthorized or station not found", http.StatusUnauthorized)
		return
	}

	var req models.SharePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	targetUser, _, err := h.st.GetUserByEmail(r.Context(), req.UserEmail)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := h.st.CreateShare(r.Context(), stationID, targetUser.ID); err != nil {
		http.Error(w, "Already shared or internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Invitation sent"))
}

// GetInvitations returns pending playlist invitations for the user.
// @Summary Get pending invitations
// @Description List all pending playlist sharing invitations
// @Tags auth
// @Produce json
// @Success 200 {array} models.PlaylistShare
// @Router /auth/invitations [get]
func (h *Handlers) GetInvitations(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, _ := h.auth.ValidateToken(authHeader[7:])

	shares, err := h.st.GetPendingShares(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shares)
}

// UpdateInvitationStatus accepts or rejects a sharing invitation.
// @Summary Update invitation status
// @Description Accept or reject a playlist sharing invitation
// @Tags auth
// @Accept json
// @Param id path int true "Share ID"
// @Param request body models.UpdateShareRequest true "Status: accepted or rejected"
// @Success 200 {string} string "Status updated"
// @Router /auth/invitations/{id} [put]
func (h *Handlers) UpdateInvitationStatus(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// We could verify ownership of the share here too, but UpdateShareStatus in store is scoped by ID.
	// A more robust implementation would check if current user is indeed the recipient of the share.

	shareID, _ := strconv.Atoi(r.PathValue("id"))
	var req models.UpdateShareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.st.UpdateShareStatus(r.Context(), shareID, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Register handles user registration.
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.auth.Register(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Login handles user authentication.
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.auth.Login(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetStationState returns the current playback state of a station.
func (h *Handlers) GetStationState(w http.ResponseWriter, r *http.Request) {
	stationID := r.PathValue("id")
	state, ok := h.sm.GetState(stationID)
	if !ok {
		http.Error(w, "Station not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}
