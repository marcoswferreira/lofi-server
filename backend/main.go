package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/marcoswferreira/lofi-server/backend/api"
	_ "github.com/marcoswferreira/lofi-server/backend/docs"
	"github.com/marcoswferreira/lofi-server/backend/models"
	"github.com/marcoswferreira/lofi-server/backend/service"
	"github.com/marcoswferreira/lofi-server/backend/store"
	"github.com/marcoswferreira/lofi-server/backend/ws"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Lofi Server API
// @version 1.0
// @description This is a lofi music streaming server with community features.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api
func main() {
	ctx := context.Background()
	
	// Load .env
	godotenv.Load()

	// Initialize Store
	st, err := store.NewStore(ctx)
	if err != nil {
		log.Printf("Could not connect to database (continuing in in-memory mode for now): %v", err)
	} else {
		if err := st.InitSchema(ctx); err != nil {
			log.Fatalf("Could not initialize schema: %v", err)
		}
	}

	sm := service.NewStationManager()
	auth := service.NewAuthService(st)

	// Initial Data (Stations)
	sm.AddStation(models.Station{
		ID:          "chill",
		Name:        "Chill Vibes",
		Description: "Smooth beats for relaxing.",
		Playlist: []models.Track{
			{ID: "jfKfPfyJRdk", Title: "lofi hip hop radio - beats to relax/study to", Duration: 24 * time.Hour},
			{ID: "5yx6BWlEVcY", Title: "Chill Lofi Mix", Duration: 10 * time.Minute},
		},
	})

	sm.AddStation(models.Station{
		ID:          "focus",
		Name:        "Focus Mode",
		Description: "Minimalist beats for deep work.",
		Playlist: []models.Track{
			{ID: "DWcJFNfaw9c", Title: "Lofi for Reading", Duration: 1 * time.Hour},
		},
	})

	h := api.NewHandlers(sm, auth, st)
	hub := ws.NewHub()
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/stations", h.GetStations)
	mux.HandleFunc("POST /api/stations", h.CreateStation)
	mux.HandleFunc("PUT /api/stations/{id}", h.UpdateStation)
	mux.HandleFunc("DELETE /api/stations/{id}", h.DeleteStation)
	mux.HandleFunc("POST /api/stations/{id}/share", h.SharePlaylist)
	mux.HandleFunc("GET /api/stations/{id}/state", h.GetStationState)
	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("GET /api/auth/invitations", h.GetInvitations)
	mux.HandleFunc("PUT /api/auth/invitations/{id}", h.UpdateInvitationStatus)
	mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	// Swagger route
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// Simple CORS wrapper
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	log.Println("Server starting on :8081")
	log.Println("Swagger UI available at http://localhost:8081/swagger/index.html")
	if err := http.ListenAndServe(":8081", corsHandler(mux)); err != nil {
		log.Fatal(err)
	}
}
