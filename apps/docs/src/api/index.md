# API Overview

The **SwampHacks API** is the core backend service that powers all technical systems, including the web dashboard and Discord bot. It acts as the central source of truth for users, applications, events, and teams.

## Purpose

The API provides structured, secure access to all hackathon data and functionality:

- Serve and validate user sessions
- Enforce role-based access and permissions
- Handle CRUD operations across domains (users, events, projects, etc.)
- Connect frontend and automation tools through a consistent interface

## Consumers

The API is designed for use by:

- **Web Dashboard**: UI for organizers, hackers, mentors, and judges
- **Discord Bot**: Handles real-time updates, commands, and integrations
- **Internal Tools**: Scripts and workflows (e.g. onboarding, analytics)

## Integration

All clients communicate with the API via HTTP using secure, authenticated requests. The system is designed for modular growth and can easily support future bots, portals, or tools.

---

This page offers a high-level overview. For deeper implementation or contributing info, see other sections of the docs.
