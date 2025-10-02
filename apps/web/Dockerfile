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
FROM caddy:2-alpine AS prod

COPY --from=build /app/dist .

COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 80
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
