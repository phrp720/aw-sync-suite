# Aw-Sync-Agent





<details>

<summary>Table of Contents</summary>

1. [About](#about)
2. [Key Features](#key-features)
3. [Requirements](#requirements)
    - [For Development](#for-development)
    - [For Running the Agent](#for-running-the-agent)
4. [Package Overview](#package-overview)
5. [Plugins](#plugins)
    - [Plugin System](#plugin-system)
    - [Available Plugins](#available-plugins)
    - [Creating Custom Plugins](#create-custom-plugins)
6. [Configuration Options](#configuration-options)
    - [Configuration Hierarchy](#configuration-hierarchy)
    - [Configurable Settings](#configurable-settings)
    - [Configuration Examples](#configuration-examples)
6. [Makefile Commands](#makefile-commands)
    - [General Commands](#general-commands)
    - [Service Commands](#service-commands)

</details>

## About

The **aw-sync-agent** is an open-source background service written in Go that  facilitates the synchronization of multiple ActivityWatch instances to a single Prometheus database, allowing for centralized monitoring and analysis of user activity data across various systems.

For a brief overview of the agent's functionality, refer to the [Agent Documentation](https://github.com/phrp720/aw-sync-suite/wiki/Agent-as-a-Service)
## Key Features

- **Data Synchronization**: Fetches user activity data from multiple ActivityWatch instances.
- **Data Filtering and Aggregation**: Filters and aggregates data based on user-defined criteria.
- **Plugin System**: Allows for easy extensibility and customization before pushing the data to Prometheus.
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

To run the agent, you need only:
- **aw-sync-agent** executable which can be built using the provided Makefile(or you can find the latest release [here](https://github.com/phrp720/aw-sync-suite/releases/)).
- **configuration file** (optional but recommended; you can also use flags or environment variables)

> [!Warning]
> - The agent requires a reachable running instance of ActivityWatch and Prometheus to function.
> - To use filters and other advanced features, you must provide the configuration file.

## Package Overview

Here is the package structure of the agent:

| Package/Folder   | Description                                                                                           |
|------------------|-------------------------------------------------------------------------------------------------------|
| **aw**           | Client for ActivityWatch REST API interactions.                                                       |
| **prometheus**   | Client for Prometheus REST API interactions.                                                          |
| **synchronizer** | Manages data synchronization from ActivityWatch to Prometheus.                                        |
| **checkpoint**   | Tracks the latest data synced for efficient operation.                                                |
| **errors**       | Error handling utilities.                                                                             |
| **datamanager**  | Handles data processing and transmission to Prometheus (**Scrape**, **Aggregate** and **Push** data). |
| **config**       | Contains the agent Configuration files.                                                               |
| **settings**     | Manages agent configuration settings.                                                                 |
| **util**         | Utility functions, including health checks.                                                           |
| **scripts**      | Additional, optional scripts.                                                                         |
| **cron**         | Manages scheduled sync intervals.                                                                     |
| **service**      | Manages service mode operations.                                                                      |
| **tests**        | Contains unit tests for the agent.                                                                    |


## Plugins

### Plugin System

The agent supports a plugin system that allows for easy extensibility and customization of the synchronization process. Plugins can be used to modify data before and after collection, apply filters, and perform other operations.

### Available Plugins

The available plugins can be found [here](https://github.com/phrp720/aw-sync-suite-plugins/wiki/%F0%9F%94%8C-Available-Plugins).

### Create Custom Plugins

To create a custom plugin, follow [these steps](https://github.com/phrp720/aw-sync-suite-plugins/wiki/%F0%9F%93%9D-How-to-Create-a-Plugin).

## Configuration Options

<p>The <code inline="">aw-sync-agent</code> offers flexibility in configuration through three layers: <strong>configuration file</strong>, <strong>environment variables</strong>, and <strong>command-line flags</strong>.

However, note that <strong>filters can only be defined in the configuration file</strong> and cannot be set via environment variables or command-line flags.

The table below details all configurable settings, their purpose, and their defaults.</p>

### Configuration Hierarchy
<p>Settings are applied in the following order of priority:</p>

1. **Command-Line Flags**: Highest priority, overrides all other settings.
2. **Environment Variables**: Override settings in the configuration file.
3. **Configuration File** (`aw-sync-settings.yaml`): Base settings, used if no overrides are provided.


### Configurable Settings

| Flag                | Environment Variable | Config Key         | Description                                                                                                                                                                           | Required | Default                  |
|---------------------|----------------------|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|--------------------------|
| -service            | -                    | -                  | Runs the agent as a continuous service.                                                                                                                                               | ❌        | -                        |
| -immediate          | -                    | -                  | Executes the synchronization process once immediately.                                                                                                                                | ❌        | -                        |
| -testConfig         | -                    | -                  | Prints the agent settings,filter(not raw) and categories results                                                                                                                      | ❌        | -                        |
| -awUrl              | ACTIVITY_WATCH_URL   | awUrl              | URL of the ActivityWatch server to fetch data from.                                                                                                                                   | ✅        | http://localhost:5600    |
| -prometheusUrl      | PROMETHEUS_URL       | prometheusUrl      | URL of the Prometheus server for sending metrics.                                                                                                                                     | ✅        | -                        |
| -prometheusAuth     | PROMETHEUS_AUTH      | prometheusAuth     | Bearer token for Prometheus authentication (useful when secured via NGINX).                                                                                                           | ❌        | -                        |
| -cron               | CRON                 | cron               | Cron expression to schedule periodic syncs .It conforms to the standard as described by the [Cron wikipedia page](https://en.wikipedia.org/wiki/Cron).                                | ❌        | Every 5 minutes          |
| -excludedWatchers   | EXCLUDED_WATCHERS    | excludedWatchers   | List of ActivityWatch watchers to exclude from syncing (use pipe `\|` for multiple entries in flags or environment variables).                                                        | ❌        | -                        |
| -plugins            | PLUGINS              | plugins            | List of [plugins](https://github.com/phrp720/aw-sync-suite-plugins) that are going to be enabled to the agent (use pipe `\|` for multiple entries in flags or environment variables). | ❌        | -                        |
| -pluginsStrictOrder | PLUGINS_STRICT_ORDER | pluginsStrictOrder | When set to true,executes the plugins with the given order                                                                                                                            | ❌        | false                    |
| -userId             | USER_ID              | userId             | Custom identifier for the user; defaults to the system hostname or a generated ID if unspecified.                                                                                     | ❌        | Hostname or Generated ID |
| -includeHostname    | INCLUDE_HOSTNAME     | includeHostname    | When set to true, appends the hostname to the exported metrics for better identification in multi-user environments.Otherwise the host value is set to `Unknown`                      | ❌        | false                    |

#### Configuration Examples
<ol>
<li>
<p><strong>Using Command-Line Flags</strong>:</p>
<pre><code class="language-bash">./aw-sync-agent -awUrl=http://localhost:5600 -prometheusUrl=http://prometheus.local -cron="*/5 * * * *"
</code></pre>
</li>
<li>
<p><strong>Using Environment Variables</strong>:</p>
<pre><code class="language-bash">export ACTIVITY_WATCH_URL=http://localhost:5600
export PROMETHEUS_URL=http://prometheus.local/api/v1/write
export CRON="*/5 * * * *"
./aw-sync-agent
</code></pre>
</li>
<li>
<p><strong>Using the Configuration File (<code inline="">aw-sync-agent.yaml</code>)</strong>:</p>
<pre><code class="language-yaml">awUrl: "http://localhost:5600"
prometheusUrl: "http://prometheus.local/api/v1/write"
cron: "*/5 * * * *"
excludedWatchers:
  - "aw-watcher-afk"
  - "aw-watcher-window"
plugins:
  - "filters"
pluginsStrictOrder: false
userId: "custom-user-id"
includeHostname: true
</code></pre>
</li>
</ol>

## Makefile Commands

### General Commands

| Command         | Description                                                           |
|-----------------|-----------------------------------------------------------------------|
| `run`           | Runs the `main.go` file.                                              |
| `build`         | Builds the `aw-sync-agent` executable.                                |
| `clean`         | Removes the `aw-sync-agent` executable.                               |
| `test`          | Runs the go tests.                                                    |
| `check-os`      | Determines the operating system by running the `detect_os.go` script. |
| `clean-service` | Cleans the ActivityWatch Sync Agent service files.                    |
| `format`        | Formats the Go code using `gofmt`.                                    |
| `build-all`     | Builds executables for both Windows and Linux.                        |
| `clean-all`     | Removes both Windows and Linux executables.                           |

### Service Commands

| Command           | Description                                                                       |
|-------------------|-----------------------------------------------------------------------------------|
| `service-install` | Builds the executable and runs it as a service.                                   |
| `service-start`   | Starts the ActivityWatch Sync Agent service.                                      |
| `service-stop`    | Stops the ActivityWatch Sync Agent service.                                       |
| `service-status`  | Displays the status of the ActivityWatch Sync Agent service.                      |
| `service-remove`  | Stops and removes the ActivityWatch Sync Agent service, and cleans service files. |
| `service-restart` | Restarts the ActivityWatch Sync Agent service.                                    |
