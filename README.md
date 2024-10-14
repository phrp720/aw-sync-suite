# aw-sync-agent

This is an open-source ActivityWatch sync agent. The purpose of this agent is to function as a service, collecting data from ActivityWatch and pushing it to a central monitoring system, Prometheus. Grafana will then be used to visualize the data.

## Test Query

```bash
curl 'http://localhost:9090/api/v1/query?query=ActivityWatch'
```

## TODO

- [ ] Create an activitywatch client to interact with ActivityWatch rest API
- [ ] Modify the already implemented prometheus client
- [ ] Create a sync agent to push data from ActivityWatch to Prometheus
- [ ] Create a Grafana dashboard to visualize the data
- [ ] Make the agent run as a service for Linux and Windows(maybe and for macOS)
- [ ] Create a docker-compose file to run the whole system.
- [ ] Create a README.md file with instructions on how to run the system