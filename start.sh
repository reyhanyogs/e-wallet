#!/bin/sh
set -e

echo "running db migration"
/app/migrate -path /app/db/migrations -database "$DB_SOURCE" -verbose up

echo "app starting"
exec "$@"