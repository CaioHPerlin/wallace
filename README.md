# WallAce

## Product

WallAce is a remotely configurable motion alarm built around an ESP32 and a PIR sensor. It records motion events, shows them in a web dashboard, and lets users control local LED and buzzer alerts remotely.

## Team

- Caio Hygino Perlin
- Caio Ricarte Loyer
- Murilo Toledo

## Project Goal

Build a small, connected monitoring system for residential, commercial, and institutional environments. WallAce provides motion detection, event history, real-time dashboard updates, and remote control of physical alerts.

## How It Works

- **Edge device:** The ESP32 reads the PIR sensor, activates the LED and buzzer according to the current configuration, and sends detections to the API over Wi-Fi.
- **Application:** The Go API stores detections and serves device configuration. The TypeScript dashboard shows event history, receives real-time updates over WebSocket, and controls the LED and buzzer settings.

## Integration for the Firmware Team

Firmware implementation and hardware assembly are assigned to **Caio Ricarte Loyer and Murilo Toledo**. The firmware sketch and flashing workflow are maintained separately; this repository contains the integration contract and local stack.

1. Start the API, dashboard, and database from the repository root with `docker compose up --build`.
2. Connect the ESP32 to the same Wi-Fi network. Set its API base URL to `http://<computer-LAN-IP>:8080`; `localhost` from the ESP32 would refer to the board itself.
3. Configure the Wi-Fi credentials and the same `API_TOKEN` used by the local stack in an untracked firmware config file. `.env.example` uses `dev-token` for local development; never commit real credentials.
4. Use the ESP32 factory MAC address as `device_id`. Implement detection reporting, configuration polling, and heartbeats according to the [OpenAPI contract](docs/api.yaml).
5. Follow [firmware integration notes](firmware/README.md) for request timing, actuator behavior, and wiring. Confirm the actual GPIO assignments during assembly and record them there.
6. Verify the full flow in the dashboard: detections should appear, and LED or buzzer changes should reach the device on its next configuration poll.

The OpenAPI document is the source of truth for endpoint paths, authentication, payloads, and dashboard messages. The [API integration guide](docs/api.md) covers local routing details.
