# Repository Guidelines

## Language

Use English for all repository content: source code, identifiers, comments, documentation, user interface text, commit messages, and pull requests. Keep conversations outside the repository in the language preferred by the people involved. Write commit messages in English using a concise Conventional Commit prefix, for example `feat: add API readiness check`.

## Project Structure

- `api/cmd/server/` is the API composition root; `api/internal/` owns configuration, persistence, HTTP handlers, and server lifecycle.
- `api/` contains the Go service and its tests.
- `web/` contains the React and TypeScript dashboard, built with Vite and served by Caddy.
- `firmware/` contains integration notes for the ESP32 team; firmware is developed and flashed separately.
- `docs/` contains the HTTP API contract and integration guidance.
- `compose.yml` runs the API, dashboard, and PostgreSQL for local development.

Keep API behavior and its contract aligned. Update `docs/api.yaml` when endpoint paths or payloads change, and link firmware contributors to that contract.

## Development Commands

- `docker compose up --build` builds and starts the local stack.
- `docker compose down` stops the stack while preserving database data.
- `cd api && go test ./...` runs Go tests; `go vet ./...` checks common issues.
- `cd web && npm install && npm run dev` starts the dashboard for frontend work.
- `cd web && npm run lint` checks TypeScript and JavaScript; `npm run format:check` checks formatting and `npm run format` applies Prettier.

## Style and Testing

Run `gofmt` on Go files. Use two spaces in TypeScript, React, JSON, and YAML; keep names descriptive and JSON fields in `snake_case`. Put Go tests beside the code in `*_test.go` files. Cover health behavior and database failures in API tests; add dashboard tests when interactive behavior is introduced.

## Pull Requests and Configuration

Describe the change, its validation, and any linked issue in each pull request. Include screenshots for visible dashboard changes. Keep credentials in an untracked `.env` file based on `.env.example`; never commit secrets. Local Compose credentials are for development only.
