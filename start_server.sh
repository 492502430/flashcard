#!/bin/bash
kill $(lsof -ti:8080) 2>/dev/null
sleep 1
export DATABASE_URL="postgres://flashcard:flash123@localhost:5433/flashcard?sslmode=disable"
/tmp/flashcard-server 2>&1 | tee /tmp/server_log.txt
