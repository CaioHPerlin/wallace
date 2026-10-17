# Firmware Integration

The firmware is maintained separately from the local dashboard and API Compose services. It targets an ESP32 using the Arduino framework, a PIR motion sensor, and local LED and buzzer actuators.

## API Integration

Use the [OpenAPI contract](../docs/api.yaml) as the source of truth for endpoint paths, payloads, and timing. It includes implemented health checks and device endpoints marked as planned. The [integration guide](../docs/api.md) explains the local routing setup. The ESP32 communicates with the API over HTTP; dashboard push uses a planned WebSocket from the API to the browser.

The local API base URL is `http://<computer-running-compose>:8080`. For an ESP32 on the same Wi-Fi network, replace the host with the computer's LAN address. Configure that address and Wi-Fi credentials outside committed source files. Authentication and retry behavior remain undecided.
