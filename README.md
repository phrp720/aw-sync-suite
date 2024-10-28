# ActivityWatch Sync Agent | aw-sync-agent
The ActivityWatch Sync Agent is a free and open-source tool designed to run as a background service. It collects data from the ActivityWatch platform and pushes it to a central monitoring system, Prometheus. This setup allows for the visualization of collected data using Grafana, providing insightful metrics on user activity.

### Features

- Data Synchronization: Retrieves and synchronizes user activity data from multiple ActivityWatch instances.
- Prometheus Integration: Transforms and pushes the data to Prometheus for storage and monitoring.
- Grafana Visualization: Easily visualize trends and activity metrics in Grafana.
- Customizable Configuration: Configure which ActivityWatch buckets to include or exclude, and set sync intervals.
- Service Mode: Run the agent as a service on Windows and Linux systems with just one command.

## Requirements

- For Development purposes:  
  - To modify the agent by yourself, you will need the following:
    - Go Version >= 1.23.2
    - Make

- For Running the agent:
  - To run the agent, you will need the following:
    - The aw-sync-agent binary
    - A running instance of ActivityWatch
    - A running instance of Prometheus
    - A running instance of Grafana(Optional)
    - A configuration file(Optional). The agent can run without a configuration file(with the use of flags or env variables), but it is recommended to have one.

### This repo contains the following packages:

- `aw`: A client to interact with the ActivityWatch REST API.
- `prometheus`: A client to interact with the Prometheus REST API.
- `synchronizer`: The synchronizer that pushes data from ActivityWatch to Prometheus.
- `checkpoint`: Contains the checkpoint mechanism to keep track of the latest data pushed.
- `errors`: Contains the error handlers.
- `datamanager`: Manages the data processing and pushing to Prometheus.
- `settings`: Handles the configuration settings for the agent.
- `util`: Contains utility functions such as health checks.
- `scripts`: Contains some useful non-mandatory scripts.
- `cron`: Contains the cron manager .
- `service`: Contains the as A Service manager.


### Settings

| Flag                | Environment Variable | Config Key          | Description                                                                               | Mandatory | Default Value                                     |
|---------------------|----------------------|---------------------|-------------------------------------------------------------------------------------------|-----------|---------------------------------------------------|
| `-service`          | -                    | -                   | Run the agent as a service.                                                               | false     | false                                             |
| `-awUrl`            | `ACTIVITY_WATCH_URL` | `aw-url`            | The URL of the ActivityWatch server.                                                      | true      | -                                                 |
| `-cron`             | `CRON`               | `cron`              | A cron expression to run the sync agent.                                                  | false     | Every 5 minutes                                   |
| `-excludedWatchers` | `EXCLUDED_WATCHERS`  | `excluded-watchers` | A pipe-separated list of watchers to exclude from the sync agent.                         | false     | -                                                 |
| `-minData`          | `MIN_DATA`           | `min-data`          | The minimum amount of data that a watcher needs to have to be included in the sync agent. | false     | 5                                                 |
| `-prometheusUrl`    | `PROMETHEUS_URL`     | `prometheus-url`    | The URL of the Prometheus server.                                                         | true      | -                                                 |
| `-userId`           | `USER_ID`            | `userId`            | The name of the user that we scrape data                                                  | false     | The name of the computer otherwise a generated id |

### Configuration Hierarchy:

1. Configuration File(config.yaml): This is the base configuration.
2. Environment Variables: These override the configuration file settings.
3. Command-Line Flags: These have the highest priority and override both the configuration file and environment variables.

### Makefile commands

- `make build`: Builds the agent.
- `make run`: Runs the agent.
- `make service`: Runs the agent as a service.
- `make format`: Formats the code.
- `make clean`: Cleans the project.


## TODO

- [x] Create an ActivityWatch client to interact with ActivityWatch rest API
- [x] Modify the already implemented prometheus client
- [x] Create a sync agent to push data from ActivityWatch to Prometheus
- [x] Create a command-line interface to run the agent(Makefile,flags,environment variables)
- [x] Create a checkpoint mechanism for optimization purposes
- [x] Create internet connection check and retry mechanism
- [x] Make the agent run as a service for Linux
- [x] Create a cron manager to run the agent at specific intervals
- [x] Create a configuration manager to handle the configuration settings(.yaml file + environment variables + flags)
- [ ] Create better error handler
- [ ] Make the agent run as a service for Windows(PENDING)
- [ ] Create an aggregator to aggregate the data
- [ ] Create a Grafana dashboard to visualize the data
- [ ] Create a Dockerfile that creates an image of the aw-sync-agent.
- [ ] Create a Documentation for the aw-sync-agent
- [ ] Publish the aw-sync-agent version 0.1.0