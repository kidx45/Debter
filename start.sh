#!/bin/sh
set -e

echo "starting migration of database"
migrate -path /debter/migration -database "$DB_URL" -verbose up

echo "starting the application"
exec "$@"