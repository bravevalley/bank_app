#! /bin/sh

set -e

echo "\033[0;33mRunning DB Migrations...\033[0m"
migrate -source /app/migrations -database "$DBSOURCE" -verbose up

echo "\033[0;32mStarting API...\033[0m"
"$1"