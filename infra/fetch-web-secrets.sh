#!/usr/bin/env bash
set -euo pipefail

# Check required environment variables
: "${INFISICAL_CLIENT_ID:?INFISICAL_CLIENT_ID is required}"
: "${INFISICAL_CLIENT_SECRET:?INFISICAL_CLIENT_SECRET is required}"
: "${WORKSPACE_SLUG:=swamphacks-core}"
: "${ENVIRONMENT:=dev}"


# Ensure jq is installed
if ! command -v jq >/dev/null 2>&1; then
  sudo apt-get update -y
  sudo apt-get install -y jq
fi

# 1) Get access token from Infisical
ACCESS_TOKEN=$(curl -sS --fail -X POST "https://us.infisical.com/api/v1/auth/universal-auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"clientId\":\"${INFISICAL_CLIENT_ID}\",\"clientSecret\":\"${INFISICAL_CLIENT_SECRET}\"}" \
  | jq -r '.accessToken // .data.accessToken')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
  echo "ERROR: Failed to obtain Infisical access token" >&2
  exit 1
fi

# 2) Fetch secrets JSON and convert to dotenv
curl -sS --fail -G "https://us.infisical.com/api/v3/secrets/raw" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-urlencode "secretPath=/web" \
  --data-urlencode "workspaceSlug=${WORKSPACE_SLUG}" \
  --data-urlencode "environment=${ENVIRONMENT}" \
  --data-urlencode "viewSecretValue=true" \
  | jq -r '.secrets[]?, .imports[].secrets[]? | "\(.secretKey)=\(.secretValue)"' \
  > ./secrets/.env.dev.web


# Optional: inspect
echo "✅ .env.dev.web generated"

