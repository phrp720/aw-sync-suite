# ActivityWatch Sync Agent | aw-sync-agent

The **ActivityWatch Sync Agent** is an open-source background service that collects data from the ActivityWatch platform and synchronizes it to a central Prometheus database. With Grafana integration, it provides real-time visual insights into user activity data, allowing for easy monitoring and analysis.

## Key Features

- **Data Synchronization**: Fetches user activity data from multiple ActivityWatch instances.
- **Prometheus Integration**: Transforms data into a Prometheus-compatible format for centralized monitoring.
- **Grafana Visualization**: Easily visualize activity metrics and trends through Grafana dashboards.
- **Flexible Configuration**: Allows selection of ActivityWatch buckets to include/exclude and customizes sync intervals.
- **Cross-Platform Service Mode**: Run as a background service on both Windows and Linux with a single command.
  - *Note*: Service mode utilizes the [service-builder](https://github.com/phrp720/service-builder) library.

## Requirements

### For Development
To modify the agent, ensure you have:
- **Go** version >= 1.23.2
- **Make**

### For Running the Agent
To run the agent, you need:
- **aw-sync-agent** binary
- Running instances of **ActivityWatch** and **Prometheus**
- **Grafana** (optional, for visualization)
- **Configuration file** (optional but recommended; you can also use flags or environment variables)

## Package Overview

- **aw**: Client for ActivityWatch REST API interactions.
- **prometheus**: Client for Prometheus REST API interactions.
- **synchronizer**: Manages data synchronization from ActivityWatch to Prometheus.
- **checkpoint**: Tracks the latest data synced for efficient operation.
- **errors**: Error handling utilities.
- **datamanager**: Handles data processing and transmission to Prometheus.
- **settings**: Manages agent configuration settings.
- **util**: Utility functions, including health checks.
- **scripts**: Additional, optional scripts.
- **cron**: Manages scheduled sync intervals.
- **service**: Manages service mode operations.

## Configuration Options

The following table provides details on configurable settings:

| Flag                | Environment Variable | Config Key          | Description                                                                               | Required | Default                                       |
|---------------------|----------------------|---------------------|-------------------------------------------------------------------------------------------|----------|-----------------------------------------------|
| `-service`          | -                    | -                   | Run the agent as a service.                                                               | No       | `false`                                       |
| `-awUrl`            | `ACTIVITY_WATCH_URL` | `aw-url`            | URL of the ActivityWatch server.                                                          | Yes      | -                                             |
| `-cron`             | `CRON`               | `cron`              | Cron expression to schedule syncs.                                                        | No       | Every 5 minutes                               |
| `-excludedWatchers` | `EXCLUDED_WATCHERS`  | `excluded-watchers` | Pipe-separated list of watchers to exclude.                                               | No       | -                                             |
| `-minData`          | `MIN_DATA`           | `min-data`          | Minimum data threshold for watchers to be included.                                       | No       | 5                                             |
| `-prometheusUrl`    | `PROMETHEUS_URL`     | `prometheus-url`    | URL of the Prometheus server.                                                             | Yes      | -                                             |
| `-userId`           | `USER_ID`            | `userId`            | Identifier for user data; defaults to computer name if not specified.                     | No       | Generated ID or computer name                 |

### Configuration Hierarchy

Settings are prioritized in the following order:
1. **Configuration File** (`config.yaml`): Base configuration settings.
2. **Environment Variables**: Override settings from the configuration file.
3. **Command-Line Flags**: Highest priority, override both file and environment settings.

## Makefile Commands

- `make build`: Builds the agent.
- `make run`: Runs the agent.
- `make service`: Starts the agent as a service.
- `make format`: Formats the codebase.
- `make clean`: Cleans up build artifacts.

## Roadmap

### Completed Tasks
- [x] ActivityWatch client for REST API interactions
- [x] Prometheus client integration
- [x] Synchronization agent to transfer data from ActivityWatch to Prometheus
- [x] CLI with Makefile, flags, and environment variable support
- [x] Checkpoint mechanism for optimized data sync
- [x] Internet connection validation and retry
- [x] Linux service mode
- [x] Cron manager for scheduled operation
- [x] Configuration management (YAML, environment variables, flags)
- [x] Windows service mode support


### Upcoming Features
- [ ] Improved error handling
- [ ] Data aggregation for enhanced insights
- [ ] Grafana dashboard template for data visualization
- [ ] Dockerfile for containerized deployment
- [ ] Complete project documentation
- [ ] Publish version 0.1.0 of aw-sync-agent
