# API Integration Guide

[`api.yaml`](api.yaml) is the canonical contract for HTTP routes and the planned dashboard WebSocket. Operation descriptions identify whether each endpoint is implemented or planned. The project-specific `x-websockets` extension describes the WebSocket message contract.

The ESP32 communicates with the API over HTTP. The dashboard will receive real-time detection events over WebSocket when that endpoint is implemented. In the local stack, Caddy forwards `/api/*` and `/ws/*` to the API service.

Authentication and retry behavior remain undecided. Keep endpoint paths, payloads, and message schemas in the OpenAPI contract rather than duplicating them here.
