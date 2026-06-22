#!/bin/bash
set -e

kill $(lsof -ti:8088) 2>/dev/null || true
sleep 1
cd "$(dirname "$0")/backend"
go build -o /tmp/flashcard-server ./cmd/server
export DATABASE_URL="postgres://flashcard:flash123@localhost:5433/flashcard?sslmode=disable"
export PORT=8088
/tmp/flashcard-server 2>&1 | tee /tmp/server_log.txt
