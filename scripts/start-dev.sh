#!/bin/bash

# 1. Load context parameters and connection secrets
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/continuum?sslmode=disable"
export NEXT_PUBLIC_API_URL="http://localhost:8080/api"

echo "🚀 Launching Continuum Backend Microservice Instance..."
cd /workspaces/Continuum/apps/api
go run cmd/server/main.go &
API_PID=$!

echo "⚛️ Starting Web Interface Server Engine..."
cd /workspaces/Continuum/apps/web
npm run dev &
WEB_PID=$!

# Handle shutdown signals clean to ensure child background routines are killed cleanly
trap "kill $API_PID $WEB_PID; exit" INT TERM EXIT

echo "🎉 Dev environment running. Press Ctrl+C to terminate all services safely."
wait