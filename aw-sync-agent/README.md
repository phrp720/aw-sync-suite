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
    - [Configurable Settings](#configurable-settings)
    - [Configuration Examples](#configuration-examples)
6. [Filters](#filters)
    - [Filter Format](#filter-format)
    - [Filter Field Descriptions](#filter-field-descriptions)
    - [Filter Example Scenario](#filter-examples)
      - [Plain Replace of Data](#plain-replace-of-data)
      - [Regex Replace of Data](#regex-replace-of-data)
      - [Drop of the Record](#drop-of-the-record)
7. [Makefile Commands](#makefile-commands)
    - [General Commands](#general-commands)
    - [Service Commands](#service-commands)

</details>

## About

The **aw-sync-agent** is an open-source background service written in Go that  facilitates the synchronization of multiple ActivityWatch instances to a single Prometheus database, allowing for centralized monitoring and analysis of user activity data across various systems.

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

To run the agent, you need only:
- **aw-sync-agent** executable which can be built using the provided Makefile(or you can find the latest release [here](https://github.com/phrp720/aw-sync-suite/releases/)).
- **configuration file** (optional but recommended; you can also use flags or environment variables)

> [!Warning]
> - The agent requires a reachable running instance of ActivityWatch and Prometheus to function.
> - To use filters and other advanced features, you must provide the configuration file.

## Package Overview

Here is the package structure of the agent:

| Package          | Description                                                                                           |
|------------------|-------------------------------------------------------------------------------------------------------|
| **aw**           | Client for ActivityWatch REST API interactions.                                                       |
| **prometheus**   | Client for Prometheus REST API interactions.                                                          |
| **synchronizer** | Manages data synchronization from ActivityWatch to Prometheus.                                        |
| **checkpoint**   | Tracks the latest data synced for efficient operation.                                                |
| **system_error** | Error handling utilities.                                                                             |
| **datamanager**  | Handles data processing and transmission to Prometheus (**Scrape**, **Aggregate** and **Push** data). |
| **settings**     | Manages agent configuration settings.                                                                 |
| **filter**       | Filters data based on user-defined criteria.                                                          |
| **util**         | Utility functions, including health checks.                                                           |
| **scripts**      | Additional, optional scripts.                                                                         |
| **cron**         | Manages scheduled sync intervals.                                                                     |
| **service**      | Manages service mode operations.                                                                      |
| **tests**        | Contains unit tests for the agent.                                                                    |

## Configuration Options

<p>The <code inline="">aw-sync-agent</code> offers flexibility in configuration through three layers: <strong>configuration file</strong>, <strong>environment variables</strong>, and <strong>command-line flags</strong>.

However, note that <strong>filters can only be defined in the configuration file</strong> and cannot be set via environment variables or command-line flags.

The table below details all configurable settings, their purpose, and their defaults.</p>

### Configuration Hierarchy
<p>Settings are applied in the following order of priority:</p>

1. **Command-Line Flags**: Highest priority, overrides all other settings.
2. **Environment Variables**: Override settings in the configuration file.
3. **Configuration File** (`aw-sync-agent.yaml`): Base settings, used if no overrides are provided.


### Configurable Settings

| Flag              | Environment Variable | Config Key       | Description                                                                                                                    | Required | Default                  |
|-------------------|----------------------|------------------|--------------------------------------------------------------------------------------------------------------------------------|----------|--------------------------|
| -service          | -                    | -                | Runs the agent as a continuous service.                                                                                        | ❌        | -                        |
| -immediate        | -                    | -                | Executes the synchronization process once immediately.                                                                         | ❌        | -                        |
| -awUrl            | ACTIVITY_WATCH_URL   | awUrl            | URL of the ActivityWatch server to fetch data from.                                                                            | ✅        | http://localhost:5600    |
| -prometheusUrl    | PROMETHEUS_URL       | prometheusUrl    | URL of the Prometheus server for sending metrics.                                                                              | ✅        | -                        |
| -prometheusAuth   | PROMETHEUS_AUTH      | prometheusAuth   | Bearer token for Prometheus authentication (useful when secured via NGINX).                                                    | ❌        | -                        |
| -cron             | CRON                 | cron             | Cron expression to schedule periodic syncs (e.g., */5 * * * * for every 5 minutes).                                            | ❌        | Every 5 minutes          |
| -excludedWatchers | EXCLUDED_WATCHERS    | excludedWatchers | List of ActivityWatch watchers to exclude from syncing (use pipe `\|` for multiple entries in flags or environment variables). | ❌        | -                        |
| -userId           | USER_ID              | userId           | Custom identifier for the user; defaults to the system hostname or a generated ID if unspecified.                              | ❌        | Hostname or Generated ID |
| -includeHostname  | INCLUDE_HOSTNAME     | includeHostname  | When set to true, appends the hostname to the exported metrics for better identification in multi-user environments.           | ❌        | false                    |

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
userId: "custom-user-id"
includeHostname: true
</code></pre>
</li>
</ol>

## Filters

### Filter Overview

The `aw-sync-agent.yaml` file supports three types of filtering, which can be applied individually or mixed for advanced filtering needs:

1. **Plain Replace**: Replaces field values with specified new values when target conditions are met.
2. **Regex Replace**: Performs partial or full replacements using regex patterns for flexible matching and substitution.
3. **Drop Record**: Removes records entirely when the specified target conditions are met.

These filtering methods can be combined in a single filter, allowing you to replace some fields while dropping others based on the same or different conditions.
### Filter Format

```yaml
Filters:

  - Filter:
    filter-name: "Plain Replace of Data" ## Name of the filter (optional)
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

    plain_replace:  ## Mapping for Values to be plain replaced
      - key: <key_name>
        value: <new_value>
        .
        .
        .
      - key: <key_name>
        value: <new_value>
  - Filter:
      filter-name: "Partial Regex Replace of data" ## Name of the filter (optional)
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

      regex_replace:  ## Mapping for Values to be replaced
        - key: <key_name>
          expression: <regex_expression>
          value: <new_value>
          .
          .
          .
        - key: <key_name>
          expression: <regex_expression>
          value: <new_value>
  - Filter:
      filter-name: "Drop of the Record" ## Name of the filter (optional)
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
    
      drop: "true" ## if true, the record will be dropped if the target conditions are met
          

```

### Filter Field Descriptions
| Field             | Description                                                                                                                                                                                                                                                                                                                                                                           |
|-------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **filter-name**   | (Optional) Specifies a name for the filter.                                                                                                                                                                                                                                                                                                                                           |
| **enabled**       | (Optional) If set to `false`, the filter will be disabled.                                                                                                                                                                                                                                                                                                                            |
| **watchers**      | (Optional) Specifies the watchers to apply the filter to, like `aw-watcher-window`. If this field is omitted or empty, the filter will apply to all watchers.                                                                                                                                                                                                                         |
| **target**        | Contains key-value pairs that the data record must match for filtering to occur. Each entry includes: <br> - **key**: The data field name, e.g., app. <br> - **value**: A regex pattern to match against the field's value.                                                                                                                                                           |
| **plain_replace** | Specifies key-value pairs for replacement. If the target key-values match, the specified keys in replace will be updated to the new values in the data record. Each entry includes: <br> - **key**: The field name to replace. <br> - **value**: The new value for the specified key.                                                                                                 |
| **regex_replace** | Specifies key-value pairs for replacement using regex patterns. If the target key-values match, the specified keys in replace will be updated to the new values in the data record. Each entry includes: <br> - **key**: The field name to replace. <br> - **expression**: A regex pattern to match against the field's value. <br> - **value**: The new value for the specified key. |
| **drop**          | If set to `true`, the record will be dropped if the target conditions are met.                                                                                                                                                                                                                                                                                                        |

### Filter Examples

#### Plain Replace of Data

This filter configuration performs a plain text replacement on data records.

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

    plain_replace:  ## key value pairs to replace e.g. on the key `title` replace its value with `Email`
      - key: "title"  ## key of record
        value: "Email" ## value to replace

```

**Explanation**:

- **watchers**: Applies this filter to `aw-watcher-window` only. If empty, the filter would apply to all watchers.
 
- **target**: Specifies matching conditions:

    - `app` must match `"Google.*"` (e.g., "Google Chrome").
    - `title` must match `"mail.*"` (e.g., "mail - Inbox").


Both conditions must match for the filter to apply.

- **plain_replace**: When the `target` conditions are met, this section replaces plain values in the matching record:
    - Sets the title field to "Email".

**Outcome**: For records in `"aw-watcher-window"` where `app` starts with "Google" and `title` starts with "mail," this filter changes the `title` field’s value to `"Email"`.
### Regex Replace of Data

This filter configuration performs a partial regex replacement on data records.

```yaml
Filters:
  - Filter:
    filter-name: "Partial Regex Replace of data" ## Name of the filter (optional)
    enable: "false" ## Enable the filter
    watchers: ## watchers where the filter will be applied (optional)
      - "aw-watcher-window"

    target: ## Data Records that if match , do the filtering (mandatory)

      - key: "app" ## key to filter on
        value: "Google.*" ## value to filter on REGEX

      - key: "title" ## key to filter on
        value: "test.*" ## value to filter on REGEX

    regex-replace: ## key value pairs to replace e.g. on the key `title` replace its value with `Email`

      - key: "title" ## key of record
        expression: "test.*" ## REGEX to replace
        value: "Email" ## value to replace
```

**Explanation**:

- **filter-name**: Specifies a name for the filter.
- **enable**: Indicates whether the filter is enabled (`false` means the filter is disabled).
- **watchers**: Applies this filter to `aw-watcher-window` only. If empty, the filter would apply to all watchers.
- **target**: Specifies matching conditions:
    - `app` must match `"Google.*"` (e.g., "Google Chrome").
    - `title` must match `"test.*"` (e.g., "test case").
- **regex-replace**: When the `target` conditions are met, this section replaces values in the matching record using regex:
    - For the `title` field, if a part of `title` matches the regex `"test.*"`, it will be replaced with `"Email"`.

### Drop of the Record

This filter configuration drops data records that match specified conditions.

```yaml
Filters:
  filter-name: "Drop of the Record" ## Name of the filter (optional)
  watchers: ## watchers where the filter will be applied (optional)
    - "aw-watcher-window"

  target: ## Data Records that if match , do the filtering (mandatory)

    - key: "app" ## key to filter on
      value: "Google.*" ## value to filter on REGEX

    - key: "title" ## key to filter on
      value: "test.*" ## value to filter on REGEX
  drop: "true" ## Drop the record if matched
```

**Explanation**:

- **filter-name**: Specifies a name for the filter.
- **watchers**: Applies this filter to `aw-watcher-window` only. If empty, the filter would apply to all watchers.
- **target**: Specifies matching conditions:
    - `app` must match `"Google.*"` (e.g., "Google Chrome").
    - `title` must match `"test.*"` (e.g., "test case").
- **drop**: If set to `true`, the record will be dropped if the `target` conditions are met.

> [!Note]
> - Filters can be combined to perform multiple operations on the same data record(plain && regex replacement).
> - Filters are applied in the order they are defined in the configuration file.
> - Filters can be disabled by setting the `enabled` field to `false`.
> - Filters that have the drop field set to `true` will not perform any replacement operations.


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