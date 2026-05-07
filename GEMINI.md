# Lofi Server - Project Documentation

## 🎵 Project Vision
A highly engaging, monetizable lofi music platform providing synchronized, multi-channel streaming combined with community features and productivity tools.

## 🏗 Architecture & Tech Stack
*   **Monorepo Structure:**
    *   `/backend`: Go (Standard library for routing, Gorilla WebSockets for real-time).
    *   `/frontend`: React + TypeScript + Vite.
*   **Database:** PostgreSQL (User profiles, Pomodoro stats, Subscription state).
*   **Real-time:** WebSockets for synchronized playback, live chat, and voting.
*   **Playback:** YouTube IFrame Player API (transparent audio playback).
*   **Monetization:** Stripe API integration.

## 🚀 Core Features
1.  **Multiple Stations:** "Chill", "Sleep", "Focus" - each with synchronized global playlists.
2.  **Productivity Suite:** Built-in Pomodoro timer with user session tracking.
3.  **Community Engagement:**
    *   Real-time chat per station.
    *   User profiles with customizable avatars/badges.
    *   Track voting (upvote/downvote) and song requests.
4.  **Monetization Streams:**
    *   **Premium Subscriptions:** Ad-free, exclusive stations, custom themes.
    *   **Virtual Goods:** Cosmetic chat items (badges, colors).
    *   **Donations/Ads:** Patreon/Kofi integration and background sponsorships.

## 🛠 Development Workflow
*   **Backend:** Port 8081. Run with `go run main.go` in `/backend`.
*   **Frontend:** Port 5173. Run with `npm run dev` in `/frontend`.
*   **State Sync:** Backend calculates `startTime` for tracks; frontend calculates `currentSeconds` to seek the YouTube player upon join.

## 📈 Status & Roadmap
- [x] Phase 1: Core Streaming & Multi-Station Setup.
- [x] Phase 2: User Accounts & Auth (Backend logic ready, requires DB).
- [x] Phase 3: Real-time Community Features (Live Chat & WS Hub).
- [x] Phase 4: Productivity Tools (Pomodoro Timer).
- [/ ] Phase 5: Monetization (Postponed).
- [x] Phase 6: Polish (Animated backgrounds, Mobile responsiveness).
- [x] Phase 7: Modernization (Ultra Premium 3D Cockpit Design, Framer Motion, 2026 UX).

## 🔒 Security & Standards
*   JWT for authentication.
*   CORS configured for local development.
*   Strict TypeScript typing in the frontend.
*   Idiomatic Go service/store pattern in the backend.
