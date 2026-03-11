# Web App Overview

The SwampHacks web app is the primary frontend for the platform. It is a single-page application that serves participants, event organizers, and platform administrators through a unified interface.

## Stack

| Component | Technology |
|---|---|
| Framework | [React 19](https://react.dev) + [Vite 6](https://vitejs.dev) |
| Language | TypeScript 5.8 |
| Package manager | npm |
| Routing | [TanStack Router](https://tanstack.com/router) |
| Data fetching | [TanStack Query](https://tanstack.com/query) |
| Forms | [TanStack Form](https://tanstack.com/form) |
| Tables | [TanStack Table](https://tanstack.com/table) |
| Accessible UI | [React Aria](https://react-spectrum.adobe.com/react-aria/) + [React Aria Components](https://react-spectrum.adobe.com/react-aria/react-aria-components.html) |
| Styling | [TailwindCSS v4](https://tailwindcss.com) |
| HTTP client | [Axios](https://axios-http.com) |
| Charts | [ECharts](https://echarts.apache.org) |
| Unit testing | [Vitest](https://vitest.dev) |
| E2E testing | [Playwright](https://playwright.dev) |
| Component explorer | [Storybook 8](https://storybook.js.org) |
| Schema validation | [Zod](https://zod.dev) |

## Feature Areas

The application is organized into feature modules under `src/features/`:

| Feature | Description |
|---|---|
| `Auth` | Discord OAuth2 login flow and session handling |
| `Onboarding` | New user profile setup after first login |
| `Dashboard` | Participant home screen with event status and quick actions |
| `Event` | Event browsing, detail views, and lifecycle state |
| `EventOverview` | Summarized event information for participants |
| `Application` | Hackathon application submission and status tracking |
| `ApplicationReview` | Reviewer interface for scoring and triaging applications |
| `FormBuilder` | Dynamic form schema builder for custom application questions |
| `CheckIn` | QR code and RFID-based attendee check-in |
| `Redeemables` | Prize and redeemable item tracking and redemption |
| `Team` | Team creation, join requests, and membership management |
| `Settings` | User account and preference settings |
| `Users` | User lookup and management utilities |
| `EventAdmin` | Organizer-facing event management controls |
| `PlatformAdmin` | Platform-wide administration (event manager, global settings) |

## Key Scripts

| Script | Command | Description |
|---|---|---|
| `dev` | `vite --host 0.0.0.0` | Start the development server |
| `build` | `vite build && tsc -b` | Production build with type checking |
| `test` | `vitest run` | Run unit tests |
| `storybook` | `storybook dev -p 6006` | Launch the Storybook component explorer |
| `generate:openapi` | `openapi-typescript ../api/docs/swagger.yaml -o ./src/lib/openapi/schema.d.ts` | Regenerate TypeScript types from the API OpenAPI spec |

## Ports

| Service | Port |
|---|---|
| Development server | `5173` |
| Storybook | `6006` |
