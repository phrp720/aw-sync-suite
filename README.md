<h1 align="center">Aw-Sync-Suite</h1>
<p align="center">
Open-Source Solution for Securely Syncing and Visualizing Multiple ActivityWatch Instances.  <br>
</p>

<p align="center">

   <a href="https://github.com/phrp720/aw-sync-suite/actions/workflows/tests.yaml?query=branch%3Amaster">
    <img title="Tests" src="https://github.com/phrp720/aw-sync-suite/actions/workflows/tests.yaml/badge.svg?branch=master" alt="tests"/>
  </a>
  <a href="https://github.com/phrp720/aw-sync-suite/actions/workflows/build.yml">
    <img title="Build Status GitHub" src="https://github.com/phrp720/aw-sync-suite/actions/workflows/build.yml/badge.svg"  alt="build"/>
  </a>
  <a href="https://github.com/phrp720/aw-sync-suite/actions/workflows/agent-docker-image.yml">
    <img title="Docker Build" src="https://github.com/phrp720/aw-sync-suite/actions/workflows/agent-docker-image.yml/badge.svg" alt="docker build">
  </a>

  <a href="https://github.com/phrp720/aw-sync-suite/releases">
    <img title="Latest release" src="https://img.shields.io/github/v/release/phrp720/aw-sync-suite" alt="Latest release">
  </a>
</p>

<p align="center">
  If you’ve ever wished for <strong> a simple, centralized solution </strong> to sync and visualize data from multiple instances of ActivityWatch, you’re in the right place.
 <br>
  📖 For detailed documentation, visit our <a href="https://github.com/phrp720/aw-sync-suite/wiki">GitHub Wiki</a>.
</p>

<details>

<summary>📑 Table of Contents</summary>

1. [About](#-about)
2. [Features](#-features)
3. [Flow Diagrams](#-flow-diagrams)
4. [Quick Start Guide](#-quick-start-guide)
    - [Download the Latest Release](#step-1-download-the-latest-release)
    - [Deploy aw-sync-center (Cloud Setup)](#step-2-deploy-aw-sync-center-cloud-setup)
    - [Configure and Run aw-sync-agent](#step-3-configure-and-run-aw-sync-agent)
    - [Visualize in Grafana](#step-4-visualize-in-grafana)
5. [Preview](#%EF%B8%8F-preview)
6. [Components](#-components)
    - [aw-sync-agent](#aw-sync-agent)
    - [aw-sync-center](#aw-sync-center)
7. [Requirements](#-requirements)
8. [Contributing](#-contributing)
</details>

## 🔍 About
**Aw-Sync-Suite** provides an easy-to-deploy solution on syncing data from multiple [ActivityWatch](https://github.com/ActivityWatch/activitywatch) instances to a centralized [Prometheus](https://prometheus.io/) database with easy visualization in [Grafana](https://grafana.com/).

The project operates independently of **ActivityWatch** and supports all ActivityWatch versions with a REST API.

### This suite consists of two main components:
- **[aw-sync-agent](https://github.com/phrp720/aw-sync-suite/tree/master/aw-sync-agent):** Runs on each device, retrieves and filters ActivityWatch data, and sends it securely to Prometheus via remote-write.
- **[aw-sync-center](https://github.com/phrp720/aw-sync-suite/tree/master/aw-sync-center):** A centralized Prometheus and Grafana setup for aggregating and visualizing data.

This repository simplifies the deployment and integration process, allowing you to monitor user activity across multiple devices with ease.

## 🌟  Features
- 🌐 **Centralized Monitoring:** Aggregate data from multiple devices effortlessly.
- 🛡️ **Data Filtering:** Protect sensitive information by filtering or sanitizing it at the source.
- 📍 **Checkpointing Mechanism:** Smart synchronization with automatic tracking of synced data.
-  📈 **Pre-Built Dashboards:** Use intuitive Grafana dashboards for instant insights.
- ⚙️ **Effortless Deployment:** Simple setup for both agent and central components.


## 📊 Flow Diagrams


| Sync-Suite Without Bearer Token Authentication 🔓 | Sync-Suite With Bearer Token Authentication 🔐          |
|---------------------------------------------------|---------------------------------------------------------|
| ![aw-sync-diagram.png](aw-sync-diagram.png)       | ![aw-sync-diagram-nginx.png](aw-sync-diagram-nginx.png) |


## 🚀 Quick Start Guide

### Step 1: Download the Latest Release

1. Visit the [Releases Page](https://github.com/phrp720/aw-sync-suite/releases/).
2. Pick the .zip file for your platform:
   - 🖥️ **Windows/Linux Agent**: Lightweight agents to sync data.
   - ☁️ **Aw-Sync-Center**: The central Prometheus-Grafana setup.
   - 📦 **Aw-Sync-Suite**: Includes everything in one bundle.
3. Extract the contents of the selected `.zip` file(s) into your desired directory.



### Step 2: Deploy **aw-sync-center** (Cloud Setup)

If you downloaded **Aw-Sync-Suite** or **Aw-Sync-Center**:

1. Navigate to the `aw-sync-center` directory:
   ```bash
   cd aw-sync-center
   ```
2. Start the cloud components (Prometheus and Grafana) using Docker Compose:
   ```bash
   docker-compose -f docker-compose-default.yaml up
   ```

This command launches all necessary services for centralized data collection and visualization.

> [!Note]
> To secure Prometheus endpoints with Bearer token authentication, follow the instructions [here](https://github.com/phrp720/aw-sync-suite/tree/master/aw-sync-center#prometheus-with-nginx-secure-setup).



### Step 3: Configure and Run **aw-sync-agent**

If you downloaded **Aw-Sync-Suite** or **Windows/Linux Agent**, follow these steps:

1. Navigate to the place where  agent is located:
    - **Aw-Sync-Suite**: `aw-sync-suite/aw-sync-agent/windows` or `aw-sync-suite/aw-sync-agent/linux`
    - **Windows Agent**: `aw-sync-agent-{version}-windows-86_64/windows`
    - **Linux Agent**: `aw-sync-agent-{version}-linux-86_64/linux`
2. Open and configure the `aw-sync-agent.yaml` file:
    - Specify the Prometheus endpoint.
    - Adjust other [settings](https://github.com/phrp720/aw-sync-suite/wiki/Agent-Configuration) and [filters](https://github.com/phrp720/aw-sync-suite/wiki/Data-Filtering) as needed.

#### Run the Agent:
You can run **aw-sync-agent** in one of the following ways:

1. **As an Executable**
    - Run the executable directly. The terminal needs to remain open:
        - Windows:
          ```cmd
          .\aw-sync-agent.exe
          ```
        - Linux:
          ```bash
          ./aw-sync-agent
          ```

2. **As a System Service**
    - Run the agent as a background service.
    - **Important**:
        - On **Windows**, you must run the terminal as an administrator to create the service successfully:
          ```cmd
          .\aw-sync-agent.exe -service
          ```
        - On **Linux**, use the following command:
          ```bash
          ./aw-sync-agent -service
          ```

3. **As a Docker Container**
    - Use Docker to run the agent in a container:
      ```bash
      docker run -v /path/to/aw-sync-agent.yaml:/opt/aw-sync-agent/aw-sync-agent.yaml phrp5/aw-sync-agent:latest
      ```
      
> [!IMPORTANT]
> Replace `/path/to/aw-sync-agent.yaml` with the actual path to your configuration file.

> [!Tip]
> - Find the latest Docker images [here](https://hub.docker.com/r/phrp5/aw-sync-agent/tags).
> - Example Docker Compose setups are available [here](https://github.com/phrp720/aw-sync-suite/tree/master/aw-sync-agent/docker-examples).
> - For detailed configuration options, check [this guide](https://github.com/phrp720/aw-sync-suite/tree/master/aw-sync-agent#configuration-options).



### Step 4: Visualize in Grafana

1. Open Grafana in your browser.
2. Add Prometheus as a data source.
3. Import the pre-built dashboards (available [here]()) to visualize ActivityWatch data.

## 👁️ Preview

| ![TOP](aw-sync-center/grafana/dashboards/screenshots/summary/summary_1.png) | ![BOTTOM](aw-sync-center/grafana/dashboards/screenshots/summary/summary_2.png) |
|-----------------------------------------------------------------------------|--------------------------------------------------------------------------------|

## 🧩 Components

### aw-sync-agent

- **Purpose**: Syncs data from ActivityWatch to Prometheus.
- **Deployment**: Run on each computer you wish to track user activity from.
- **Configuration**: Configure it via the `aw-sync-agent.yaml` file.

### aw-sync-center

- **Purpose**: Centralized cloud setup that includes Prometheus and Grafana for monitoring and visualization.
- **Deployment**: Set up once for centralized control and management.
- **Included Services**: Prometheus, Grafana, and necessary dashboards.

## 🛠️ Requirements

- Docker and Docker Compose for easy setup of `aw-sync-center`.
- A running instance of ActivityWatch on the computers you want to monitor.

## 👥 Contributing
Contributions are welcomed! If you have ideas, improvements, or bug fixes, feel free to open an issue or submit a pull request.

