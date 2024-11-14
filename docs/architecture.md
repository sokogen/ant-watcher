# Architecture Overview

This document provides an overview of the architecture of the **ant-watcher** service.

## High-Level Architecture

The **ant-watcher** service is designed to collect and aggregate metrics from GitHub Actions workflows and jobs. It processes incoming webhook events, interacts with the GitHub API, and sends formatted metrics to monitoring systems like VictoriaMetrics or Prometheus.

### Components

- **Main Application (`cmd/ant-watcher/`)**

  The entry point of the application. It initializes the configuration, sets up logging, and starts the necessary components such as servers and workers.

- **Reusable Packages (`pkg/`)**

  - **`pkg/githubapi`**

    Handles interactions with the GitHub API, including fetching workflow and job information, handling rate limits, and caching responses.

  - **`pkg/metrics`**

    Formats and sends metrics to the monitoring system. Supports both `/metrics` exposition for Prometheus and `/push` endpoint for VictoriaMetrics.

  - **`pkg/memory`**

    Manages in-memory objects, implements TTL (Time-To-Live), and handles cleanup of old objects to prevent memory overconsumption.

  - **`pkg/config`**

    Loads and validates configuration from files, environment variables, and command-line flags.

  - **`pkg/logger`**

    Provides structured logging with support for different log levels (info, warning, error) and formats logs in JSON for easy parsing.

  - **`pkg/errors`**

    Defines custom error types and provides centralized error handling mechanisms.

- **Internal Packages (`internal/`)**

  - **`internal/webhook`**

    Handles incoming webhook events from GitHub, validates and authenticates requests, and enqueues tasks for asynchronous processing.

  - **`internal/collector`**

    Processes data from webhooks and the GitHub API, aggregates metrics, and prepares them for sending to the monitoring system.

  - **`internal/middleware`**

    Implements middleware functions for HTTP servers, such as authentication, logging, rate limiting, and IP whitelisting.

- **API Handlers (`api/`)**

  Manages HTTP endpoints for health checks, metrics exposition, and any additional API functionalities.

- **Configuration (`config/`)**

  Contains configuration files, such as `config.json`, and handles parsing and validation of configuration parameters.

- **Scripts (`scripts/`)**

  Includes scripts for building, testing, and deploying the application, as well as other automation tasks.

## Data Flow

1. **Service Initialization**

   - The main application loads the configuration and initializes logging.
   - HTTP servers for webhooks and API endpoints are started.
   - Workers and queues for asynchronous processing are initialized.

2. **Webhook Processing**

   - The service receives a webhook event from GitHub.
   - The request is validated and authenticated using the webhook secret and allowed IPs.
   - Valid events are enqueued for asynchronous processing.

3. **Asynchronous Processing**

   - Workers dequeue tasks and process the events.
   - The service interacts with the GitHub API to fetch additional data if necessary.
   - Metrics are generated and stored in memory with TTL.

4. **Metrics Dispatching**

   - Metrics are exposed via the `/metrics` endpoint for Prometheus scraping.
   - If configured, metrics are also pushed to VictoriaMetrics via the `/push` endpoint.
   - After successful dispatch, metrics are removed from memory.

5. **Memory Management**

   - The service monitors memory usage.
   - If memory usage exceeds configured limits, old metrics are purged based on TTL or using a priority queue.
   - Garbage collection metrics are monitored to optimize performance.

6. **Graceful Shutdown**

   - The service handles termination signals.
   - Ongoing processes are completed.
   - Resources are cleaned up before exit.

## Considerations

- **Scalability**

  While the service is designed to handle a moderate load, it can be scaled horizontally. Stateless components and externalizing state can improve scalability.

- **Error Handling**

  Centralized error handling and retries are implemented to improve resilience against transient failures.

- **Security**

  - Webhook requests are validated using secrets and IP whitelisting.
  - Sensitive configurations (like tokens and secrets) are managed securely.
  - Logs are sanitized to prevent leakage of sensitive information.

## Future Enhancements

- **Distributed Tracing**

  Integration with OpenTelemetry for detailed performance analysis.

- **Advanced Rate Limiting**

  Improved handling of GitHub API rate limits and dynamic adjustments.

- **Configuration Management**

  Support for dynamic reloading of configuration without restarting the service.

---
