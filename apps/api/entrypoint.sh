#!/bin/bash

# Define colors
BOLD="\033[1m"
RESET="\033[0m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
CYAN="\033[0;36m"
RED="\033[0;31m"

# Spacer
divider() {
  echo -e "${CYAN}\n------------------------------------------------------------\n${RESET}"
}

echo -e "${BOLD}${GREEN}🚀 Bootstrapping Dev Environment${RESET}"

divider

# Load environment variables
if [ -f .env.local ]; then
  echo -e "${YELLOW}📦 Loading environment variables from .env.local...${RESET}"
  export $(cat .env.local | xargs)
else
  echo -e "${RED}⚠️  .env.local not found. Skipping env load.${RESET}"
fi

divider

# Run database migrations
echo -e "${YELLOW}📂 Running database migrations...${RESET}"
make migrate-up

divider

# Generate sqlc code
echo -e "${YELLOW}🛠️  Generating sqlc code...${RESET}"
make generate

divider

# Start the app with hot reload
echo -e "${GREEN}🔥 Starting application with Air (hot reload)...${RESET}"
exec air
