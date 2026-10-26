# Firmware Integration

The firmware is maintained separately from the API and dashboard Compose services. It targets an ESP32 using the Arduino framework, a PIR motion sensor, and local LED and buzzer actuators. This document is the integration contract; the sketch itself is developed and flashed separately.

The [OpenAPI contract](../docs/api.yaml) is the source of truth for endpoint paths, payloads, and timing. The [integration guide](../docs/api.md) explains local routing and identity.

## Network

The local API base URL is `http://<computer-running-compose>:8080`. For an ESP32 on the same Wi-Fi network, replace the host with the computer's LAN address. Configure the base URL and Wi-Fi credentials outside committed source files (for example, in a local `secrets.h` that is gitignored).

The ESP32 speaks plain HTTP to the API. TLS is not used in the local stack.

## Device identity

A device is identified by its ESP32 factory MAC address, used verbatim as `device_id` (for example `A4:CF:12:34:56:78`). Read it on the station interface with `WiFi.macAddress()` after `WiFi.begin()`. No registration is required: the API creates the device row on first contact.

The dashboard shows a user-facing `name` that defaults to the MAC until renamed.

## Authentication

Every request under `/api/devices` must send a shared bearer token:

```
Authorization: Bearer <token>
```

The token matches the API's `API_TOKEN`. Store it in the same untracked config as the URL and Wi-Fi credentials.

## Message flow and timing

| Action | Endpoint | Method | Payload | When |
| --- | --- | --- | --- | --- |
| Poll config | `/api/devices/{mac}/config` | GET | — | Every 3–5 seconds |
| Heartbeat | `/api/devices/{mac}/heartbeat` | POST | — | Every 15–30 seconds |
| Report detection | `/api/devices/{mac}/detections` | POST | `{"ts": <unix seconds>}` | On PIR trigger |
| Read own status | `/api/devices/{mac}/status` | GET | — | Optional |

Keep the heartbeat interval under 60 seconds: the API marks a device offline when its last heartbeat is older than 60 seconds.

The detection payload `ts` is the device clock's Unix timestamp in seconds. Send it promptly after the PIR triggers; if the device has no real-time clock, a monotonic millisecond counter plus a boot epoch is enough for the dashboard to order events.

## Actuator config

`GET /config` returns:

```json
{ "led_enabled": true, "buzzer_enabled": false }
```

The firmware should apply this on every poll: light the LED only while `led_enabled` is true, and sound the buzzer on a detection only while `buzzer_enabled` is true. The dashboard updates these values; the firmware never writes config.

## Wiring

Pins are placeholders until the prototype hardware is chosen. Confirm and record the real assignments before assembly.

| Signal | Type | Connection |
| --- | --- | --- |
| PIR output | digital input | GPIO input (3.3 V logic, often with an onboard regulator) |
| LED | digital output | GPIO through a current-limiting resistor |
| Buzzer | digital output | GPIO through a transistor or current-limiting resistor |

## Suggested loop

1. Connect Wi-Fi, read the MAC, start an HTTP client with the shared token.
2. Poll `GET /config` and apply `led_enabled`/`buzzer_enabled`.
3. Watch the PIR input; on a rising edge, `POST /detections` with the current timestamp and trigger the buzzer if enabled.
4. Every 15–30 seconds, `POST /heartbeat`.
5. Back off briefly on HTTP failures and retry; the API is reachable once Compose is healthy.

Retry and backoff behavior remains undecided; keep it simple for the prototype (a short fixed delay is sufficient).
