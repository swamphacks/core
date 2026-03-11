# Welcome to SwampHacks Core

Welcome to the official developer documentation for the **SwampHacks Core** platform. 

This platform serves as the central nervous system for the SwampHacks event. It handles hacker registration, team formation, application review, project submissions, and day-of-event operations. 

Whether you are a new organizer joining the technical team or looking to understand how the platform operates under the hood, this documentation will guide you through our architecture, local setup, and deployment workflows.

---

## The Tech Stack

Our platform is built with a modern, scalable stack designed for rapid development and high availability during the event:

* **Frontend:** React + Vite
* **Backend:** Go 
* **Database:** Neon DB (Serverless Postgres)
* **Secrets Management:** Infisical
* **Infrastructure:** Docker & DigitalOcean
* **Community:** Discord Bot (Custom)

---

## How to Use These Docs

Use the top navigation tabs to explore different domains of the platform:

* **Home:** Start here! Understand the high-level architecture, how the repository is structured, and how to get your local development environment running.
* **Web:** Everything related to the frontend user interface, state management, and API integration.
* **API:** Deep dive into the backend business logic, endpoints, database schema, and authentication.
* **Discord Bot:** Documentation for our custom Discord bot, including commands and event handling.
* **Infrastructure:** The DevOps playbook. Learn how we use Docker, manage secrets with Infisical, and deploy to DigitalOcean.
* **Operations:** Crucial information for event day, troubleshooting guides, and credentials for third-party services.

!!! tip "Just getting started?"
    If you are setting up the project for the first time, head straight over to the [Getting Started](getting-started.md) guide to spin up your local environment.

---

## Contributing

As an organizer, you are encouraged to keep these docs updated. If you add a new API endpoint, change the database schema, or update a deployment script, please update the corresponding page in these docs so the knowledge is passed down to the next team!
