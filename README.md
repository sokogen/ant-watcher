# ant-watcher

**Author**: Gennadii Sokolachko <g.sokolachko@moonactive.com>

## Description

`ant-watcher` is a service for collecting and aggregating metrics from GitHub Actions workflows and jobs, then sending them to monitoring systems like VictoriaMetrics or Prometheus.

## Project Structure

- `cmd/ant-watcher/` — Main application entry point.
- `pkg/` — Reusable packages (external systems).
- `internal/` — Internal packages used only within the service.
- `config/` — Configuration files.
- `scripts/` — Scripts for automation tasks.
- `docs/` — Documentation.

## Other documentation

- [Installation Guide](docs/installation.md)
- [Architecture Overview](docs/architecture.md)
- [API Reference](docs/api.md)
