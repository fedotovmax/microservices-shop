#!/bin/sh
set -e

echo "Running database migrations..."
/app/migrator -m up

echo "Start user service"
exec /app/user_service "$@"