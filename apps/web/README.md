# SwampHacks Event Management Portal

### Overview

This is the official event management portal for SwampHacks, built with modern web technologies including **Vite**, **React**, and **TypeScript**. It leverages:

- 🧭 [TanStack Router](https://tanstack.com/router/latest) for routing
- 🔄 [TanStack Query](https://tanstack.com/query/latest) + [Axios](https://axios-http.com/) for data fetching and caching
- 🧩 [React Aria](https://react-spectrum.adobe.com/react-aria/index.html) for accessible and customizable UI components

### Setup Instructions

To get started:

1. Clone the repo and make sure [pNPm](https://pnpm.io/) is installed on your system.

```bash
git clone https://github.com/swamphacks/core.git
cd core/apps/web
```

2. Install dependencies

```bash
pnpm install
```

3. Configure environment variables:

```bash
cp .env.example .env
```

Fill in the required keys and tokens in your new `.env` file. The prefix `SWAMPHACKS` is required for the environment variables to load properly.

4. Finally, launch the app

```bash
pnpm run dev
```
