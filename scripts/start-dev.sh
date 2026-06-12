#!/bin/bash

# Source the central environment variables into this bash context shell
if [ -f /workspaces/Continuum/.env ]; then
    export $(cat /workspaces/Continuum/.env | grep -v '^#' | xargs)
fi

echo "🚀 Launching Continuum Backend Microservice Instance..."
cd /workspaces/Continuum/apps/api
go run cmd/server/main.go &
API_PID=$!

echo "⚛️ Starting Web Interface Server Engine..."
cd /workspaces/Continuum/apps/web
npm run dev &
WEB_PID=$!

trap "kill $API_PID $WEB_PID; exit" INT TERM EXIT

echo "🎉 Dev environment running. Press Ctrl+C to terminate all services safely."
wait