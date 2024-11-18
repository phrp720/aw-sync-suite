# Aw-Sync-Agent





<details>

<summary>Table of Contents</summary>

1. [About](#about)
2. [Key Features](#key-features)
3. [Requirements](#requirements)
    - [For Development](#for-development)
    - [For Running the Agent](#for-running-the-agent)
4. [Package Overview](#package-overview)
5. [Configuration Options](#configuration-options)
    - [Configuration Hierarchy](#configuration-hierarchy)
6. [Filters](#filters)
    - [Filter Format](#filter-format)
    - [Filter Field Descriptions](#filter-field-descriptions)
    - [Filter Example Scenario](#filter-example)
7. [Makefile Commands](#makefile-commands)
8. [Roadmap](#roadmap)
    - [Completed Tasks](#completed-tasks)
    - [Upcoming Features](#upcoming-features)
</details>

## About

The **aw-sync-agent** is an open-source background service written in Go that collects data from [ActivityWatch](https://github.com/ActivityWatch/activitywatch) and synchronizes it to a central Prometheus database. With Grafana integration, it provides real-time visual insights into user activity data, enabling easy monitoring and analysis.

The **aw-sync-agent** facilitates the synchronization of multiple ActivityWatch instances to a single Prometheus database, allowing for centralized monitoring and analysis of user activity data across various systems.


This project operates independently of [ActivityWatch](https://github.com/ActivityWatch/activitywatch) and is compatible with all versions of ActivityWatch that support the REST API feature.

## Key Features

- **Data Synchronization**: Fetches user activity data from multiple ActivityWatch instances.
- **Data Filtering and Aggregation**: Filters and aggregates data based on user-defined criteria.
- **Prometheus Integration**: Transforms data into a Prometheus-compatible format for centralized monitoring.
- **Grafana Visualization**: Easily visualize activity metrics and trends through Grafana dashboards.
- **Flexible Configuration**: Allows selection of ActivityWatch buckets to include/exclude and customizes sync intervals.
- **Cross-Platform Service Mode**: Run as a background service on both Windows and Linux with a single command.
    - *Note*: Service mode utilizes the [go-service-builder](https://github.com/phrp720/go-service-builder) library.

## Requirements

### For Development
To modify the agent, ensure you have:
- **Go** version >= 1.23.2
- **Make**

### For Running the Agent
To run the agent, you need:
- **aw-sync-agent** executable
- **configuration file** (optional but recommended; you can also use flags or environment variables)
- Running instances of **ActivityWatch** and **Prometheus**
- Running instance of **Grafana** (optional, for visualization)

## Package Overview

- **aw**: Client for ActivityWatch REST API interactions.
- **prometheus**: Client for Prometheus REST API interactions.
- **synchronizer**: Manages data synchronization from ActivityWatch to Prometheus.
- **checkpoint**: Tracks the latest data synced for efficient operation.
- **system_error**: Error handling utilities.
- **datamanager**: Handles data processing and transmission to Prometheus(**Scrape**,**Aggregate** and **Push** data).
- **settings**: Manages agent configuration settings.
- **filter**: Filters data based on user-defined criteria.
- **util**: Utility functions, including health checks.
- **scripts**: Additional, optional scripts.
- **cron**: Manages scheduled sync intervals.
- **service**: Manages service mode operations.

## Configuration Options

The following table provides details on configurable settings:

| Flag                | Environment Variable | Config Key          | Description                                                          | Required | Default                  |
|---------------------|----------------------|---------------------|----------------------------------------------------------------------|----------|--------------------------|
| `-service`          | -                    | -                   | Runs the agent as a service.                                         | ❌        | -                        |
| `-immediate`        | -                    | -                   | Runs the synchronizer once immediately.                              | ❌        | -                        |
| `-awUrl`            | `ACTIVITY_WATCH_URL` | `aw-url`            | URL of the ActivityWatch server.                                     | ✅        | http://localhost:5600    |
| `-prometheusUrl`    | `PROMETHEUS_URL`     | `prometheus-url`    | URL of the Prometheus server.                                        | ✅        | -                        |
| `-prometheusAuth`   | `PROMETHEUS_AUTH`    | `prometheus-auth`   | Bearer Auth for prometheus(if prom is protected)                     | ❌        | -                        |
| `-cron`             | `CRON`               | `cron`              | Cron expression to schedule syncs.                                   | ❌        | Every 5 minutes          |
| `-excludedWatchers` | `EXCLUDED_WATCHERS`  | `excluded-watchers` | List of watchers to exclude(Pipe-separated for env or flag).         | ❌        | -                        |
| `-userId`           | `USER_ID`            | `userId`            | Identifier for user nickname; defaults to hostname if not specified. | ❌        | hostname or Generated ID |

### Configuration Hierarchy

Settings are prioritized in the following order:
1. **Configuration File** (`aw-sync-agent.yaml`): Base configuration settings and filtering.
2. **Environment Variables**: Override settings from the configuration file.
3. **Command-Line Flags**: Highest priority, override both file and environment settings.


## Filters

This guide explains the rules for configuring filters in the `aw-sync-agent.yaml` file, allowing you to filter data records based on key-value conditions and replace values as specified.

### Filter Format

```yaml
Filters:

  - Filter:
    watchers: ##(Optional) watchers where the filter will be applied. If empty, the filter will apply to all watchers
      - <watcher_name>

    target:  ## Conditions that if match , it will apply the filtering for the specific record
      - key: <key_name>
        value: <value_to_match> 
        .
        .
        .
      - key: <key_name>
        value: <value_to_match>

    replace:  ## Mapping for Values to be replaced
      - key: <key_name>
        value: <new_value>
        .
        .
        .
      - key: <key_name>
        value: <new_value>
      

```

### Filter Field Descriptions

- **watchers**(Optional): Specifies the watchers to apply the filter to, like `aw-watcher-window`. If this field is omitted or empty, the filter will apply to all watchers.

- **target**: Contains key-value pairs that the data record must match for filtering to occur. Each entry includes:
    - **key**: The data field name, e.g., app.
    - **value**: A regex pattern to match against the field's value.

- **replace**: Specifies key-value pairs for replacement. If the target key-values match, the specified keys in replace will be updated to the new values in the data record. Each entry includes:
    - **key**: The field name to replace.
    - **value**: The new value for the specified key.

### Filter Example

```yaml
Filters:

  - Filter:
    watchers:  ## watchers where the filter will be applied. If empty, the filter will apply to all watchers
      - "aw-watcher-window"

    target: ## Data Records that if match , do the filtering for the specific record
      - key: "app" ## key to filter on
        value: "Google.*" ## value to filter on RegEX
      
      - key: "title" ## key to filter on     
        value: "mail.*"  ## value to filter on RegEX

    replace:  ## key value pairs to replace e.g. on the key `title` replace its value with `Email`
      - key: "title"  ## key of record
        value: "Email" ## value to replace

```

**Explanation**:

- **watchers**: Applies this filter to `aw-watcher-window` only. If empty, the filter would apply to all watchers.

- **target**: Specifies matching conditions:

    - `app` must match `"Google.*"` (e.g., "Google Chrome").
    - `title` must match `"mail.*"` (e.g., "mail - Inbox").


Both conditions must match for the filter to apply.

- **replace**: When the `target` conditions are met, this section replaces values in the matching record:
    - Sets the title field to "Email".

**Outcome**: For records in `"aw-watcher-window"` where `app` starts with "Google" and `title` starts with "mail," this filter changes the `title` field’s value to `"Email"`.


## Makefile Commands

- `make build`: Builds the agent.
- `make run`: Runs the agent.
- `make service-install`: Install and Starts the agent as a service.
- `make format`: Formats the codebase.