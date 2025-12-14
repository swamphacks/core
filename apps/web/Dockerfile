# ========== Base Stage ===========
FROM node:22.16.0-slim AS base

WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm

COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

COPY . .

# ========== Dev Stage ===========
FROM base AS dev

EXPOSE 5173

CMD ["pnpm", "run", "dev", "--", "--host", "0.0.0.0"]

# ========== Build Stage ===========
FROM base AS build

# Build the React app (output will be in /app/dist)
RUN pnpm run build

# ========== Production Stage (With Runtime Config Injection) ===========
FROM node:22.16.0-slim AS prod
WORKDIR /app

# Copy built React app (dist folder)
COPY --from=build /app/dist ./dist

# Install 'serve' to serve static files
RUN npm install -g serve

# --- Runtime Config Setup ---

# Copy the entrypoint script into the working directory
# NOTE: You must ensure entrypoint.sh is in the same directory as your Dockerfile
COPY entrypoint.sh .

# Make the script executable
RUN chmod +x entrypoint.sh

# Set the entrypoint to run the script first
ENTRYPOINT ["./entrypoint.sh"]

# --- Serving Setup ---

# Expose port for top-level Caddy to proxy
EXPOSE 80

# This CMD will be passed as arguments ("$@") to the entrypoint script
CMD ["serve", "-s", "dist", "-l", "80"]
