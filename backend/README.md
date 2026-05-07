# Lofi Server - Backend 🚀

This is the backend for the Lofi Server project, built with Go. It handles synchronized music streaming, user authentication, real-time chat via WebSockets, and productivity tracking.

## 🛠 Tech Stack
- **Language:** Go 1.25+
- **Routing:** Standard `net/http` library.
- **WebSockets:** [Gorilla WebSocket](https://github.com/gorilla/websocket).
- **Database:** PostgreSQL with [pgx](https://github.com/jackc/pgx).
- **Authentication:** JWT (JSON Web Tokens) via `golang-jwt/jwt`.
- **API Documentation:** Swagger/OpenAPI.

## 📂 Project Structure
- `/api`: HTTP handlers and routing logic.
- `/models`: Data structures for Stations, Users, etc.
- `/service`: Business logic (Auth, Station management).
- `/store`: Database access layer.
- `/ws`: WebSocket hub and connection management.
- `/docs`: Swagger documentation.

## 🚀 Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL database

### Environment Variables
Create a `.env` file in the `/backend` directory (or use environment variables):
- `DATABASE_URL`: `postgres://user:password@localhost:5432/lofi_server`
- `JWT_SECRET`: Your secret key for JWT signing.
- `PORT`: Server port (default: 8081).

### Running Locally
1. Install dependencies:
   ```bash
   go mod download
   ```
2. Run the server:
   ```bash
   go run main.go
   ```

### API Documentation
The API is documented using Swagger. You can find the documentation in the `/docs` folder. To update documentation, use `swag init`.

## 🔌 Core Logic: State Sync
The backend calculates the `startTime` for tracks in each station's playlist. When a user joins, the frontend receives the current track and its start time, allowing it to calculate the `currentSeconds` to seek the YouTube player and achieve global synchronization.
