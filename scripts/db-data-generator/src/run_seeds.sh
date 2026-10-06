#!/usr/bin/env bash
#
# Runs all seed scripts in dependency order:
#   users, workshops -> no dependencies
#   hackathons       -> no dependencies
#   redeemables      -> depends on hackathons (FK on hackathon_id)
#
# Usage:
#   Put DATABASE_URL=postgresql://user:password@host:port/dbname in a
#   .env file next to this script, then:
#     ./run_seeds.sh
#   or set it inline / export it and run:
#     export DATABASE_URL="postgresql://user:password@host:port/dbname"
#     ./run_seeds.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTHON="${PYTHON:-python3}"

# Load variables from a .env file next to this script, if present.
ENV_FILE="${SCRIPT_DIR}/.env"
if [ -f "$ENV_FILE" ]; then
    set -o allexport
    source "$ENV_FILE"
    set +o allexport
fi

if [ -z "${DATABASE_URL:-}" ]; then
    echo "Error: DATABASE_URL is not set." >&2
    echo "Set it in a .env file next to this script, or run:" >&2
    echo "  export DATABASE_URL=\"postgresql://user:password@host:port/dbname\"" >&2
    exit 1
fi

run_seed() {
    local script="$1"
    echo "==> Running ${script}"
    "$PYTHON" "${SCRIPT_DIR}/${script}"
    echo "==> Done: ${script}"
    echo
}

# Order matters: hackathons before redeemables (FK dependency).
run_seed "generate_users.py"
run_seed "generate_workshops.py"
run_seed "generate_hackathon.py"
run_seed "generate_redeemables.py"

echo "All seed scripts completed."