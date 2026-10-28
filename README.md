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

## Hardware Components and Wiring

### Required components

- ESP32 DevKit and a USB cable for power and firmware flashing
- HC-SR501 PIR motion sensor
- One LED and a 220–330 Ω series resistor
- Active buzzer for simple on/off alerts, or a passive buzzer for tone/PWM output
- Breadboard and jumper wires
- For a bare buzzer: an NPN transistor and a base resistor; add a flyback diode only if it is an inductive buzzer. A driver module can be used instead.

### Suggested wiring

The following pin map is a starting point for a classic ESP32-WROOM DevKit. Verify the labels on the team's exact board and record the final assignments in [`firmware/README.md`](firmware/README.md).

| Component pin           | Connect to                                  |
| ----------------------- | ------------------------------------------- |
| HC-SR501 `VCC`          | ESP32 `5V` or `VIN`                         |
| HC-SR501 `GND`          | ESP32 `GND`                                 |
| HC-SR501 `OUT`          | ESP32 `GPIO27`                              |
| LED anode (long leg)    | ESP32 `GPIO18` through a 220–330 Ω resistor |
| LED cathode (short leg) | ESP32 `GND`                                 |

For a bare buzzer, use `GPIO19` to control an NPN transistor through a base resistor (1 kΩ is a starting point for a small buzzer). Connect the transistor emitter to `GND`, its collector to the buzzer's negative pin, and the buzzer's positive pin to the voltage marked on the buzzer. Connect the ESP32, PIR, and buzzer supply grounds together. Check the transistor's pinout before wiring. For an inductive buzzer, add a flyback diode across the buzzer, with its cathode toward the positive supply. For a driver module, connect its `SIG`/`IN`, supply, and `GND` pins according to its labels. A passive buzzer needs a tone/PWM signal; an active buzzer can be switched on and off.

The common HC-SR501 specification lists a 4.8–20 V supply and a 3.3 V output. Module variants can differ, so verify `OUT` before connecting it: ESP32 GPIOs must not receive a 5 V signal. See the [HC-SR501 specification](https://www.datocms-assets.com/28969/1662716328-hc-sr501.pdf), [Espressif DevKitC pinout](https://docs.espressif.com/projects/esp-dev-kits/en/latest/esp32/esp32-devkitc/user_guide.html), and [Espressif GPIO voltage guidance](https://docs.espressif.com/projects/esp-faq/en/latest/esp-faq-en-master.pdf).

## Integration for the Firmware Team

Firmware implementation and hardware assembly are assigned to **Caio Ricarte Loyer and Murilo Toledo**. The firmware sketch and flashing workflow are maintained separately; this repository contains the integration contract and local stack.

1. Start the API, dashboard, and database from the repository root with `docker compose up --build`.
2. Connect the ESP32 to the same Wi-Fi network. Set its API base URL to `http://<computer-LAN-IP>:8080`; `localhost` from the ESP32 would refer to the board itself.
3. Configure the Wi-Fi credentials and the same `API_TOKEN` used by the local stack in an untracked firmware config file. `.env.example` uses `dev-token` for local development; never commit real credentials.
4. Use the ESP32 factory MAC address as `device_id`. Implement detection reporting, configuration polling, and heartbeats according to the [OpenAPI contract](docs/api.yaml).
5. Follow [firmware integration notes](firmware/README.md) for request timing and actuator behavior. Use the wiring above as a starting point, confirm the final GPIO assignments during assembly, and record them there.
6. Verify the full flow in the dashboard: detections should appear, and LED or buzzer changes should reach the device on its next configuration poll.

The OpenAPI document is the source of truth for endpoint paths, authentication, payloads, and dashboard messages. The [API integration guide](docs/api.md) covers local routing details.
