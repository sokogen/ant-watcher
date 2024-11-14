# API Endpoints

This document describes the API endpoints provided by the **ant-watcher** service.

## Table of Contents

- [Webhook Endpoint](#webhook-endpoint)
- [Metrics Endpoint](#metrics-endpoint)
- [Health Check Endpoint](#health-check-endpoint)
- [Push Metrics Endpoint](#push-metrics-endpoint)

## Webhook Endpoint

### `POST /webhook`

- **Description**: Receives webhook events from GitHub.
- **Headers**:
  - `X-Hub-Signature-256`: The HMAC hex digest of the request body, used for validating the payload.
- **Body**: The JSON payload of the webhook event.
- **Authentication**: Validated using the webhook secret configured in `config.json`.
- **IP Whitelisting**: Only accepts requests from allowed IPs specified in the configuration.
- **Responses**:
  - `200 OK`: The event was received and enqueued for processing.
  - `400 Bad Request`: Invalid request or failed validation.
  - `401 Unauthorized`: Authentication failed.
  - `403 Forbidden`: IP address not allowed.
  - `500 Internal Server Error`: An error occurred while processing the request.

## Metrics Endpoint

### `GET /metrics`

- **Description**: Exposes metrics in Prometheus format for scraping.
- **Authentication**: None (ensure network security or add authentication if needed).
- **Responses**:
  - `200 OK`: Returns the metrics in text format.
  - `500 Internal Server Error`: An error occurred while generating metrics.

## Health Check Endpoint

### `GET /health`

- **Description**: Provides a simple health check for the service.
- **Authentication**: None.
- **Responses**:
  - `200 OK`: The service is running properly.
  - `500 Internal Server Error`: The service is experiencing issues.

## Push Metrics Endpoint

### `POST /push`

- **Description**: Receives metrics data to be pushed to the monitoring system.
- **Authentication**: Configurable; can require authentication tokens.
- **Body**: Metrics data in the appropriate format.
- **Responses**:
  - `200 OK`: Metrics were successfully received and processed.
  - `400 Bad Request`: Invalid metrics data.
  - `401 Unauthorized`: Authentication failed.
  - `500 Internal Server Error`: An error occurred while processing metrics.

---

**Note**: The `/push` endpoint is optional and used when pushing metrics to systems like VictoriaMetrics that support metric ingestion via HTTP POST requests.

---
