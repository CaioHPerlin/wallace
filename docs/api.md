# API Integration

The initial API implements only `/api/health/live` and `/api/health/ready`. The device endpoints are planned and documented in [`api.yaml`](api.yaml); firmware and dashboard integrations should treat them as the target contract until they are implemented.

The API uses HTTP between the ESP32 and server. The dashboard will receive real-time detection events through `ws://<host>/ws/dashboard` after that WebSocket endpoint is implemented. Each event will use `{"type":"detection","device_id":"...","ts":1732300000}`, with `ts` as Unix time in seconds. Caddy forwards both `/api/*` and `/ws/*` to the API service.

Device detection timestamps (`ts`) are Unix timestamps in seconds. The ESP32 should poll device configuration every 3–5 seconds and send a heartbeat every 15–30 seconds. Authentication and retry behavior remain to be decided.
