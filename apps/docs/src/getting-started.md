# Getting Started

Welcome to the monorepo for our hackathon event management system. This repository contains all major services powering our platform, including:

- **Web** – The frontend dashboard for organizers and attendees.
- **API** – The core backend service responsible for business logic, data management, and integrations.
- **Discord Bot** – A companion bot for community engagement and operations. (Currently independent, may depend on the API in the future.)

Each project has its own dedicated documentation tab for detailed setup and usage. This page covers the global setup needed to get the monorepo running locally.

---

## Prerequisites

Ensure the following are installed:

- [Docker Engine](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/) (if not included with Docker)
- [Git](https://git-scm.com)
- `make` (optional, but helpful for managing workflows)

---

## Initial Setup

1. **Clone the repository**:
``` bash title="Terminal/Shell"
git clone https://github.com/swamphacks/core

cd core
```
2. **Set Up the API**  
The API powers the platform's business logic, database operations, and integrations.  
Please follow the [API's Installation Guide](api/installation.md) to complete the setup.
Once you're done, return here to continue with the other services.

3. **Set Up the Web Dashboard**  
After the API is running, you can set up the frontend dashboard for organizers and attendees.  
Follow the [Web Installation Guide](web/installation.md) to install dependencies, configure environment variables, and run the app locally.


4. **Set Up the Discord Bot (OPTIONAL)**  
The Discord Bot helps manage community engagement and provides real-time updates during the event.  
Refer to the [Discord Bot Installation Guide](discord-bot/installation.md) for instructions on setup, permissions, and development workflow.

Now that you are done with setting up each project, continue down below to developing using docker.

---

## Running Development with Docker

!!! warning "Ensure Running In Root Directory"
    The following commands to begin development with docker **MUST** be run in root. Your working directory should be `core/`.

Now that you have all the main projects set up, you can use docker and docker compose to quickly start a development environment.

### **Docker Compose**  
We can use `docker compose` to quickly start up both our API and Web projects as well as local databases and caches.
In the root of the monorepo, run the following in your terminal:  
``` bash title="Terminal/Shell"
docker compose up
```
### **Docker Compose with Rebuilding**
Sometimes we need to rebuild the docker environment and container in order to apply changes. This is more of a rare instance but in case you do need to ensure everything is rebuilt, run docker with the `--build` flag.  
``` bash title="Terminal/Shell"
docker compose up --build
```
### **Docker Compose for API Development**
For all our backend developers out there, most of the time you don't need to spin up every application in order to use and test the backend. In cases where you only want to run the API without the web, specify `api` in the command.  
``` bash title="Terminal/Shell"
docker compose up api
```

!!! info "Database and Cache"
    Although it won't start the web service anymore, the database and redis cache **WILL** be started every single time the API service runs as it depends on the database and redis cache to function.

---

> ✅ Be sure to check that the API is running properly before launching the Web or Bot services, as they may rely on API endpoints.
