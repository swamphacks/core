# ========== Base Stage ===========
FROM node:22.16.0-slim AS base

WORKDIR /app

RUN npm install -g pnpm

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

# ========== Dev Stage ===========
FROM base AS dev

EXPOSE 5173

CMD ["pnpm", "run", "dev", "--", "--host", "0.0.0.0"]

# ========== Build Stage ===========
FROM base AS build
RUN pnpm run build

# ========== Production Stage ===========
FROM node:22.16.0-slim AS prod
WORKDIR /app

# Copy built React app
COPY --from=build /app/dist ./dist

# Install serve to serve static files
RUN npm install -g serve

# Expose port for top-level Caddy to proxy
EXPOSE 80

# Start serve
CMD ["serve", "-s", "dist", "-l", "80"]
