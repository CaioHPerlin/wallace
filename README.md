# WallAce

WallAce is a remotely configurable motion alarm built around an ESP32 and a PIR sensor. The dashboard shows service availability today; event history, remote actuator controls, and real-time updates are planned for later implementation.

## Local Stack

Requirements: Docker with the Compose plugin. Copy `.env.example` to `.env` to customize local PostgreSQL credentials, then start the stack:

```sh
docker compose up --build
```

Open the dashboard at <http://localhost:5173>. The API is available at <http://localhost:8080/api/health/ready>. PostgreSQL runs on the Compose network and is not published to the host.

```sh
docker compose logs -f
docker compose down
```

Compose keeps database files in the `postgres_data` volume. To remove the stack and its local database contents, run `docker compose down -v`.

## Working Locally

Run the database with `docker compose up -d db`. Start the API from `api/` with `DATABASE_URL='postgres://wallace:wallace-dev-password@localhost:5432/wallace?sslmode=disable' go run ./cmd/server`. Start the dashboard from `web/` with `npm install && npm run dev`; Vite forwards `/api` requests to the local API.

The API exposes `/api/health/live` and `/api/health/ready`; the latter checks PostgreSQL connectivity. See [`docs/api.yaml`](docs/api.yaml) for the canonical API contract, including planned device routes and dashboard WebSocket messages, and [`firmware/README.md`](firmware/README.md) for ESP32 integration notes.

## Language and Contributions

Repository content and commit messages are in English. Follow [AGENTS.md](AGENTS.md) for development commands, formatting, testing, and pull request guidance.
