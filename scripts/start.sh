#!/bin/sh
# =============================================================================
# Debter — Container entrypoint script
# =============================================================================
# This script is the Docker ENTRYPOINT. It:
#   1. Waits for the database to be reachable (avoids race conditions with
#      Docker Compose depends_on + healthcheck)
#   2. Runs database migrations via golang-migrate
#   3. Drops privileges to appuser and execs the main binary
#
# The exec + su-exec pattern replaces the shell process entirely, so PID 1
# is the application (not the shell), which ensures proper signal handling.
# =============================================================================

set -e

# ---------------------------------------------------------------------------
# Phase 1: Wait for database readiness
# ---------------------------------------------------------------------------
# Even though Docker Compose depends_on + healthcheck should ensure the DB
# is ready, this is a safety net. pg_isready is a lightweight check that
# returns 0 when the database accepts connections.
echo "[entrypoint] Checking if database health already tested ..."
if [ "$IS_CHECK_DONE" != "true"]; then
    echo "[entrypoint] Waiting for database ..."
    until pg_isready -d $DB_URL > /dev/null 2>&1; do
        echo "[entrypoint] Database not ready, retrying in 1s..."
        sleep 1
    done
fi
echo "[entrypoint] Database is ready."

# ---------------------------------------------------------------------------
# Phase 2: Run database migrations
# ---------------------------------------------------------------------------
# golang-migrate applies SQL files from the migration directory in order.
# Using -verbose so migration progress is visible in container logs.
echo "[entrypoint] Running database migrations..."
./migrate -path /debter/migration -database "$DB_URL" -verbose up
echo "[entrypoint] Migrations complete."

# ---------------------------------------------------------------------------
# Phase 3: Execute the application as non-root user
# ---------------------------------------------------------------------------
# su-exec replaces the current process (running as root) with the command
# run as appuser. This means the Go binary runs without root privileges,
# limiting the impact of any container escape or vulnerability.
echo "[entrypoint] Starting application as appuser..."
exec su-exec appuser "$@"
