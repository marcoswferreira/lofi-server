# Lofi Server 🎵

Lofi Server is a highly engaging lofi music platform providing synchronized, multi-channel streaming combined with community features and productivity tools.

## 🚀 Vision
A platform that brings people together through lofi music, offering synchronized playback across different stations ("Chill", "Sleep", "Focus"), real-time chat, and a built-in Pomodoro timer for productivity.

## 🏗 Architecture & Tech Stack

The project is structured as a monorepo:

- **Backend:** [Go](./backend) (Standard library for routing, Gorilla WebSockets for real-time).
- **Frontend:** [React + TypeScript + Vite](./frontend) (Framer Motion for animations, YouTube IFrame Player API for playback).
- **Database:** PostgreSQL (User profiles, Pomodoro stats, Subscription state).
- **Real-time:** WebSockets for synchronized playback, live chat, and voting.

## 🛠 Features

1.  **Multiple Stations:** "Chill", "Sleep", "Focus" with synchronized global playlists.
2.  **Productivity Suite:** Built-in Pomodoro timer with user session tracking.
3.  **Community Engagement:** Real-time chat per station and user profiles.
4.  **Modern UI:** Ultra Premium 3D Cockpit Design with smooth animations.

## 🚦 Quick Start

### Using Docker Compose (Recommended)

The easiest way to get the entire stack running is using Docker Compose:

```bash
docker-compose up --build
```

- **Frontend:** [http://localhost](http://localhost) (Port 80 in Docker, 5173 for local dev)
- **Backend:** [http://localhost:8081](http://localhost:8081)
- **Database:** PostgreSQL on port 5432

### Local Development

#### Prerequisites
- Go 1.25+
- Node.js 20+
- PostgreSQL 15+

#### Backend
1. Navigate to `/backend`
2. Run `go run main.go`

#### Frontend
1. Navigate to `/frontend`
2. Install dependencies: `npm install`
3. Run development server: `npm run dev`

## 🔒 Security
- JWT for authentication.
- Environment variables for sensitive data (see `docker-compose.yml`).

## 📈 Roadmap
- [x] Phase 1: Core Streaming & Multi-Station Setup.
- [x] Phase 2: User Accounts & Auth.
- [x] Phase 3: Real-time Community Features (Live Chat).
- [x] Phase 4: Productivity Tools (Pomodoro Timer).
- [ ] Phase 5: Monetization (Stripe Integration).
- [x] Phase 6: Polish & Mobile Responsiveness.
- [x] Phase 7: Modernization (2026 UX).

---
Developed as part of the Lofi Server project.
