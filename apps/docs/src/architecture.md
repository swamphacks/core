# System Architecture

The SwampHacks Core platform is designed around a central API that acts as the "brain" of the operation, coordinating data between our web client, background workers, external bots, and databases.

## 1. The Core API (The Brain)
* **Tech:** Golang
* **Role:** Acts as the central source of truth. It exposes RESTful endpoints, enforces authentication and authorization, and directly manages all database transactions. All other services (frontend, bot, workers) interact with or are driven by this API.

## 2. Background Workers
To keep the Core API performant, heavy or asynchronous tasks are offloaded to dedicated background workers written in Go.

* **Queue Management:** Redis is used to manage the task queues. (Note: Redis is currently strictly for worker queues but is provisioned to handle API caching in the future).

* **Email Worker:** Listens for events and handles reliable, asynchronous email dispatch to users.

* **BAT Worker (Balanced Admissions Thresher):** Executes the complex logic for hacker admissions. While currently implemented as a worker, its logic is decoupled enough to be extracted into a standalone service if scaling requires it.

## 3. Web Frontend
* **Tech:** React, TanStack Router

* **Role:** A lightweight single-page application (SPA) providing the user interface for hackers, judges, and organizers. It handles client-side routing via TanStack Router and relies entirely on the Core API for state and data persistence.

## 4. Discord Bot
* **Tech:** Python (discord.py)

* **Role:** Manages the SwampHacks Discord server. It interfaces directly with the Core API to perform tasks such as linking hacker accounts to Discord profiles, managing role assignments based on platform status, and triggering announcements.

## 5. Data Layer
* **Primary Database:** Neon (Serverless PostgreSQL). Stores all persistent application data, including users, applications, teams, and scoring metrics.

* **Message Broker / Cache:** Redis.

## 6. Infrastructure & Deployment
* **Containerization:** Docker is used universally. Both local development and production environments run identical Docker containers to eliminate "works on my machine" issues.

* **CI/CD:** GitHub Actions manages our pipelines, automatically running linters and triggering deployments upon merged pull requests.

* **Secrets Management:** Infisical. Hardcoded secrets do not exist in the repository. During deployment, the production server pulls the necessary environment variables directly from Infisical before spinning up the containers.
