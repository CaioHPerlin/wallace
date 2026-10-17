# Firmware Integration

The firmware is maintained separately from the local dashboard and API Compose services. It targets an ESP32 using the Arduino framework, a PIR motion sensor, and local LED and buzzer actuators.

## Planned Device API

See the [OpenAPI contract](../docs/api.yaml) and [integration notes](../docs/api.md). The API endpoints are planned and are not implemented in the current stack. The current transport decision is HTTP from the ESP32 to the API; dashboard push will use WebSocket from the API to the browser.

- `POST /api/devices/{device_id}/detections` when motion is detected; `ts` is Unix time in seconds.
- `GET /api/devices/{device_id}/config` every 3–5 seconds to apply LED and buzzer settings locally.
- `POST /api/devices/{device_id}/heartbeat` every 15–30 seconds.

The local API base URL is `http://<computer-running-compose>:8080`. For an ESP32 on the same Wi-Fi network, replace the host with the computer's LAN address. Configure that address and Wi-Fi credentials outside committed source files. Authentication and retry behavior are open decisions.
