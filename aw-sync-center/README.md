# Aw-Sync-Center

The **aw-sync-center** is a **centralized monitoring and reporting solution** using **Grafana and Prometheus**, designed to collect and visualize user activity data.

## Overview

This repository provides a **Docker Compose setup** to deploy Grafana and Prometheus as a centralized hub, collecting data from ActivityWatch Sync Agents ([aw-sync-agent](https://github.com/phrp720/aw-sync-agent)) monitoring user activity on various systems.

## Getting Started

### Prerequisites
- **Docker** and **Docker Compose** installed on your machine.

### Setup Options

This repository offers two Docker Compose configurations:
- **docker-compose-with-nginx**: Includes NGINX with Bearer token authentication for added security.
- **docker-compose-default**: A simpler setup without authentication.

### Installation and Running

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/aw-sync-center.git
   cd aw-sync-center

2. **Choose a Docker Compose Configuration**:

   - For the secure setup with Bearer token authentication:
    ```bash
    docker-compose -f docker-compose-with-nginx.yml up -d
    ```
   - For the default setup without authentication:
     ```bash
     docker-compose -f docker-compose-default.yml up -d
      ```
3. **Access Grafana**:

   - Open a browser and go to http://localhost:3000 to access Grafana.
     - Default login credentials are admin:admin .
## Prometheus with NGINX (Secure Setup)

The **docker-compose-with-nginx** configuration uses an **NGINX reverse proxy** to protect Prometheus endpoints (`/api/v1/write` and `/-/healthy`) with Bearer token authentication.

#### Generating Bearer Tokens

To generate tokens for authentication, use the `createBearerToken.py` script. This will create a `tokens.conf` file in the NGINX directory with the specified tokens, allowing secure access to Prometheus.

1. Run the following command:
```bash
python3 createBearerToken.py
```
2. **Follow the prompts** to specify the number of tokens. The script will output a new `tokens.conf` file(if not exists) within the nginx directory.
3. Using the Generated Tokens
   - Requests sent to Prometheus endpoints through NGINX must include a valid token.
   - Each token will be checked against `tokens.conf` for authentication.


### Project Status

#### In Progress
- Docker Compose setup for Grafana and Prometheus
- Prometheus authentication configuration

#### Upcoming Features
- Grafana dashboards for data visualization
- Expanded documentation for setting up the Sync Center
