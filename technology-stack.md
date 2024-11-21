# Technology Stack

| Technology     | Description                                                                     |
|----------------|---------------------------------------------------------------------------------|
| Go             | The main programming language used for developing the application.              |
| Makefile       | Used for automating the build process, running tests, and other tasks.          |
| GitHub Actions | Used for CI/CD pipelines.                                                       |
| ActivityWatch  | Used for tracking the user's activity.                                          |
| Prometheus     | Used as the central data store, with data pushed via remote-write mode.         |
| Grafana        | Used for visualizing the metrics collected by Prometheus.                       |
| Docker         | Used for containerizing Prometheus, Grafana, and the agent for easy deployment. |