#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Source the central environment variables into this bash context shell
if [ -f "$ROOT_DIR/.env" ]; then
    set -a
    # shellcheck disable=SC1091
    source "$ROOT_DIR/.env"
    set +a
fi

echo "🚀 Launching Continuum Backend Microservice Instance..."
cd "$ROOT_DIR/apps/api"
go run cmd/server/main.go &
API_PID=$!

echo "⚛️ Starting Web Interface Server Engine..."
cd "$ROOT_DIR/apps/web"
npm run dev &
WEB_PID=$!

trap "kill $API_PID $WEB_PID; exit" INT TERM EXIT

echo "🎉 Dev environment running. Press Ctrl+C to terminate all services safely."
wait
