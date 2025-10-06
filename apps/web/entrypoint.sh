#!/bin/sh

# The location where the built static files (including the placeholder env-config.js)
# are copied in the production stage.
CONFIG_FILE_PATH="/app/dist/env-config.js"

echo "Generating runtime configuration file at $CONFIG_FILE_PATH..."

# Safely inject environment variables into the config file
cat <<EOF > $CONFIG_FILE_PATH
window.ENV = {
  VITE_BASE_API_URL: "${VITE_BASE_API_URL}",
  VITE_DISCORD_OAUTH_CLIENT_ID: "${VITE_DISCORD_OAUTH_CLIENT_ID}",
  VITE_ALLOWED_HOSTS: '${VITE_ALLOWED_HOSTS}'
};
EOF

echo "Injecting BASE_API_URL = ${VITE_BASE_API_URL}"
echo "Injecting discord oauth id = ${VITE_DISCORD_OAUTH_CLIENT_ID}"
echo "Injecting allowed hosts = ${VITE_ALLOWED_HOSTS}"

echo "Runtime configuration generation complete."

# Execute the main command passed to the container (the 'serve' command from CMD)
exec "$@"
