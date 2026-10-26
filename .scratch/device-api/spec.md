# Device API, Firmware Handoff, and Dashboard

Status: resolved

## Goal

Make the device endpoints real so the ESP32 firmware team can implement and test against a live API, complete the firmware integration notes, and wire the dashboard to real data.

## Decisions

- **Device identity**: `device_id` is the ESP32 factory MAC address (e.g. `A4:CF:12:...`). The API auto-creates the device row on first contact. A user-editable `name` field makes devices identifiable in the dashboard; empty name falls back to the MAC.
- **Auth**: shared bearer token (`Authorization: Bearer <token>`) required on all `/api/devices/**` routes. Health and the dashboard WebSocket stay open. The token comes from `API_TOKEN`; the firmware keeps it in untracked config, the dashboard receives it at build time via `VITE_API_TOKEN`.
- **Firmware artifact**: integration notes only (no sketch). Firmware is developed and flashed separately.

## Contract changes (`docs/api.yaml`)

- Add `name` field and a `Device` schema.
- Add `GET /api/devices` (list) and `PATCH /api/devices/{device_id}` (rename).
- Add bearer-token security scheme on device routes.
- Mark device routes and the WebSocket as implemented.

## Phases

1. **Contract** — update `docs/api.yaml` and `docs/api.md`.
2. **API** — schema bootstrap in `database`, device store methods, device handlers, bearer middleware, WebSocket hub, wire in `router.go`/`main.go`, config gains `API_TOKEN`.
3. **Firmware notes** — expand `firmware/README.md` (message flow, timings, MAC identity, token/URL config, wiring).
4. **Front** — device list, rename, status, config toggles, detection history, real-time WebSocket.
5. **Handoff checklist** — verify the ESP32 team has everything.

## Auth default

- `Authorization: Bearer <token>` header (chosen default).
- Dashboard routes are protected by the same token and get it via `VITE_API_TOKEN`; WebSocket `/ws/dashboard` is open (push-only).

## Validation

- `cd api && gofmt -l . && go vet ./... && go test ./...`
- `cd web && npm install && npm run build && npm run lint && npm run format:check`

## Outcome

All phases implemented and validated. The device store lives in `api/internal/database/device.go` (schema bootstrap in `postgres.go`), device handlers plus bearer middleware in `api/internal/http/device.go`, and the WebSocket hub in `ws.go`. The dashboard renders devices, rename, config toggles, detection history, and real-time events. Firmware integration notes cover identity (MAC), auth, timing, and wiring.
