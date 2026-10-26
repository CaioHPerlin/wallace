# API Integration Guide

[`api.yaml`](api.yaml) is the canonical contract for HTTP routes and the dashboard WebSocket. Operation descriptions identify whether each endpoint is implemented or planned. The project-specific `x-websockets` extension describes the WebSocket message contract.

The ESP32 communicates with the API over HTTP. The dashboard receives real-time detection events over the dashboard WebSocket. In the local stack, Caddy forwards `/api/*` and `/ws/*` to the API service.

## Device identity

A device is identified by its ESP32 factory MAC address (e.g. `A4:CF:12:34:56:78`), used verbatim as `device_id`. The API creates a device on first contact, so no registration is required. The dashboard can set a user-facing `name`; when empty, the MAC is displayed.

## Authentication

Device routes under `/api/devices` require a shared bearer token: `Authorization: Bearer <token>`. The token is configured with `API_TOKEN` (API) and the matching value in the firmware's untracked config and the dashboard's `VITE_API_TOKEN`. Health routes and the dashboard WebSocket are open.

Retry behavior remains undecided. Keep endpoint paths, payloads, and message schemas in the OpenAPI contract rather than duplicating them here.
