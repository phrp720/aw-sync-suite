# Aw-Sync-Center

<details>
<summary>Table of Contents</summary>

1. [About](#about)
2. [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Setup Options](#setup-options)
    - [Installation and Running](#installation-and-running)
3. [Secure Prometheus with NGINX](#secure-prometheus-with-nginx)
    - [Generating Bearer Tokens](#generating-bearer-tokens)
4. [Folder Structure](#folder-structure)

</details>


## About

**Aw-Sync-Center** is a centralized monitoring and visualization solution built with **Prometheus** and **Grafana**. It seamlessly collects and displays user activity data synchronized by ActivityWatch Sync Agents ([aw-sync-agent](https://github.com/phrp720/aw-sync-suite/tree/master/aw-sync-agent)). The setup supports both secure and simple configurations to meet various deployment needs.



## Getting Started

### Prerequisites

Before proceeding, ensure the following are installed on your system:

- **Docker**
- **Docker Compose**

### Setup Options

The repository provides two Docker Compose configurations:

1. **Secure Setup (Recommended):**  
   Includes an NGINX reverse proxy to secure Prometheus endpoints using Bearer token authentication.  
   File: `docker-compose-with-nginx.yml`

2. **Default Setup:**  
   A simple configuration without authentication.  
   File: `docker-compose-default.yml`


### Installation and Running

1. **Download the Aw-Sync-Center**:  
   Download the `.zip` file from the [aw-sync-suite releases](https://github.com/phrp720/aw-sync-suite/releases) and extract it, or clone the repository:
   ```bash
   git clone https://github.com/phrp720/aw-sync-suite.git
   cd aw-sync-suite/aw-sync-center
   ```

2. **Choose Your Setup**:
    - **Default Setup (No Authentication):**
      ```bash
      docker-compose -f docker-compose-default.yml up -d
      ```

    - **Secure Setup (With NGINX Authentication):**
        - First, generate Bearer tokens (refer to [Generating Bearer Tokens](#generating-bearer-tokens)).
        - Then run:
          ```bash
          docker-compose -f docker-compose-with-nginx.yml up -d
          ```

3. **Access Grafana**:
    - Open [http://localhost:3000](http://localhost:3000) in your browser.
    - Default credentials:
        - Username: `admin`
        - Password: `admin`

4. **Import Dashboards**:
    - Navigate to the Grafana UI and upload dashboards from the `grafana/dashboards` directory.
    - The dashboards are pre-configured for user activity data.


## Secure Prometheus with NGINX

The secure configuration includes an NGINX reverse proxy to protect critical Prometheus endpoints (`/api/v1/write` and `/-/healthy`) using Bearer token authentication. Only authorized clients with valid tokens can access these endpoints.

### Generating Bearer Tokens

To create Bearer tokens:

1. Run the token generation script:
   ```bash
   python3 createBearerToken.py
   ```

2. Follow the prompts to specify how many tokens you need.
    - The script will generate or update a `tokens.conf` file located in the `nginx` directory.

3. Include the token in requests to authenticate with Prometheus:
   ```bash
   curl -H "Authorization: Bearer <your-token>" http://localhost:8080/api/v1/write
   ```



## Folder Structure

| Directory    | Description                                                        |
|--------------|--------------------------------------------------------------------|
| `grafana`    | Contains pre-configured Grafana dashboards and data sources.       |
| `prometheus` | Holds Prometheus configuration files.                              |
| `nginx`      | NGINX configuration for secure setups, including token management. |



