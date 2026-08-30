# =============================================================================
# Debter — Multi-stage Dockerfile
# =============================================================================
# Stage 1 (builder): Compiles the Go binary
# Stage 2 (runtime): Minimal alpine image with only runtime dependencies
#
# Environment variables are injected at container start via Docker Compose
# env_file directive — NOT baked into the image. This keeps secrets out of
# image layers and out of version control.
# =============================================================================

# ---------------------------------------------------------------------------
# Stage 1: Builder — compile Go binary
# ---------------------------------------------------------------------------
FROM golang:1.27.0-alpine3.24 AS builder

WORKDIR /debter
COPY . .
RUN go build -o main cmd/main.go
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


# ---------------------------------------------------------------------------
# Stage 2: Runtime — minimal alpine with only what we need
# ---------------------------------------------------------------------------
FROM alpine:3.24 AS runtime

# Install runtime dependencies:
#   su-exec          — drop privileges from root to appuser (required for ENTRYPOINT)
#   postgresql-client — provides pg_isready for database health checks in start.sh
RUN apk add --no-cache su-exec postgresql-client

# Create non-root user group for the application.
# Running as root inside a container is a security risk: a container escape
# or vulnerability in the app would grant root on the host. Using appuser
# limits the blast radius.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /debter

# Copy only the artifacts needed at runtime — not the full source tree.
# This keeps the image small (~15MB vs ~400MB with full Go toolchain).
COPY --from=builder /debter/main .
COPY --from=builder /debter/internal/migration ./migration
COPY --from=builder /debter/scripts/start.sh ./scripts/start.sh
COPY --from=builder /go/bin/migrate .
RUN chmod +x scripts/start.sh

# Application will listen on this port
EXPOSE 8080

# scripts/start.sh runs migrations, then su-exec drops privileges to appuser
# and runs the main binary. This ensures:
#   1. Migrations run with sufficient privileges (if needed)
#   2. The application itself runs as a non-root user
ENTRYPOINT ["./scripts/start.sh"]
CMD ["./main"]
