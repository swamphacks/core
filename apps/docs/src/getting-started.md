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
```
2. **Set Up the API**  
The API powers the platform's business logic, database operations, and integrations.  
Please follow the [API's Getting Started Guide](api/getting-started.md) to complete the setup.
Once you're done, return here to continue with the other services.

3. **Set Up the Web Dashboard**  
After the API is running, you can set up the frontend dashboard for organizers and attendees.  
Follow the [Web Getting Started Guide](web/getting-started.md) to install dependencies, configure environment variables, and run the app locally.


4. **Set Up the Discord Bot**  
The Discord Bot helps manage community engagement and provides real-time updates during the event.  
Refer to the [Discord Bot Getting Started Guide](bot/getting-started.md) for instructions on setup, permissions, and development workflow.

Now that you are done with setting up each project, continue down below to developing using docker.

---

## Running Development with Docker

Now that you have all the main projects set up, you can use docker and docker compose to quickly start a development environment.

1. **Docker Compose**  
We can use `docker compose` to quickly start up both our API and Web projects as well as local databases and caches.
In the root of the monorepo, run the following in your terminal:  
```bash "Terminal"
docker compose -f docker-compose.dev.yml
```


---

> ✅ Be sure to check that the API is running properly before launching the Web or Bot services, as they may rely on API endpoints.
