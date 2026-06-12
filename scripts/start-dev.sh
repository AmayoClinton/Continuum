#!/bin/bash

#!/bin/bash

# Resolve repository root (script lives in ./scripts)
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# Source the central environment variables into this bash context shell
if [ -f "$REPO_ROOT/.env" ]; then
    export $(cat "$REPO_ROOT/.env" | grep -v '^#' | xargs)
fi

# Ensure a writable Go module cache is available (avoid permission issues)
GOMODCACHE=${GOMODCACHE:-"$REPO_ROOT/.cache/go-mod"}
mkdir -p "$GOMODCACHE"
export GOMODCACHE

echo "🚀 Launching Continuum Backend Microservice Instance..."
cd "$REPO_ROOT/apps/api" || exit 1
go run cmd/server/main.go &
API_PID=$!

echo "⚛️ Starting Web Interface Server Engine..."
cd "$REPO_ROOT/apps/web" || exit 1
if ! [ -d node_modules ]; then
    echo "📦 Installing web dependencies..."
    npm ci --no-audit --no-fund
fi

# start Next dev on port 3000 to avoid colliding with backend
PORT=3000 npm run dev &
WEB_PID=$!

trap "kill $API_PID $WEB_PID; exit" INT TERM EXIT

echo "🎉 Dev environment running. Press Ctrl+C to terminate all services safely."
wait