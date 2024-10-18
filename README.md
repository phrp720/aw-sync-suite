# aw-sync-agent

This is an open-source ActivityWatch sync agent. The purpose of this agent is to function as a service, collecting data from ActivityWatch and pushing it to a central monitoring system, Prometheus. Grafana will then be used to visualize the data.

## Test Query

```bash
curl 'http://localhost:9090/api/v1/query?query=aw_watcher_window'
```

### Settings

| Flag                | Environment Variable | Description                                                                               | Mandatory | Default                                           |
|---------------------|----------------------|-------------------------------------------------------------------------------------------|-----------|---------------------------------------------------|
| `-asService`        | `AS_SERVICE`         | Run the agent as a service.                                                               | false     | false                                             |
| `-awUrl`            | `ACTIVITY_WATCH_URL` | The URL of the ActivityWatch server.                                                      | true      | -                                                 |
| `-cron`             | `CRON`               | A cron expression to run the sync agent.                                                  | false     | Every 10 minutes                                  |
| `-excludedWatchers` | `EXCLUDED_WATCHERS`  | A comma-separated list of watchers to exclude from the sync agent.                        | false     | -                                                 |
| `-minData`          | `MIN_DATA`           | The minimum amount of data that a watcher needs to have to be included in the sync agent. | false     | 5                                                 |
| `-prometheusUrl`    | `PROMETHEUS_URL`     | The URL of the Prometheus server.                                                         | true      | -                                                 |
| `-userID`           | `USER_ID`            | The name of the user that we scrape data                                                  | false     | The name of the computer otherwise a generated id |


### What we expect:

    sudo ./agent -excludedWatchers=aw-watcher-window -cron=*2*** -minData=9 -asService=true -awUrl=http://localhost:5600 -prometheusUrl=http://localhost:9090 -userID=Phillip

## TODO

- [x] Create an activitywatch client to interact with ActivityWatch rest API
- [x] Modify the already implemented prometheus client
- [x] Create a sync agent to push data from ActivityWatch to Prometheus
- [x] Create a command-line interface to run the agent
- [ ] Create checkpoints with checkpoint.json file
- [ ] Create internet connection check and retry mechanism
- [ ] Create better error handler
- [ ] Create a Grafana dashboard to visualize the data
- [ ] Make the agent run as a service for Linux and Windows(maybe and for macOS)
- [ ] Create a docker-compose file to run the whole system.
- [ ] Create a README.md file with instructions on how to run the system