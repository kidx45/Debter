#!/usr/bin/env bash
# =============================================================================
# Debter — Deploy to local Minikube cluster
# =============================================================================
# 1. Verifies minikube and kubectl are installed; exits early if not.
# 2. Starts Minikube (idempotent — continues if already running).
# 3. Creates the debter-config ConfigMap and debter-secrets Secret from
#    the .env.config and .env.dev.pg files (unquoted so --from-env-file
#    reads clean values).
# 4. Applies every manifest in the k8s/ directory.
#
# Usage: ./scripts/deploy.sh
# =============================================================================

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# ---------------------------------------------------------------------------
# 1. Prerequisites
# ---------------------------------------------------------------------------
command -v minikube >/dev/null 2>&1 \
  || { echo "[deploy] ERROR: minikube is not installed or not on PATH. Exiting." >&2; exit 1; }
command -v kubectl >/dev/null 2>&1 \
  || { echo "[deploy] ERROR: kubectl is not installed or not on PATH. Exiting." >&2; exit 1; }

echo "[deploy] Prerequisites OK (minikube + kubectl found)."

# ---------------------------------------------------------------------------
# 2. Start Minikube (continues if already running)
# ---------------------------------------------------------------------------
echo "[deploy] Starting minikube..."
minikube start || echo "[deploy] minikube already running or start returned non-zero; continuing."

# ---------------------------------------------------------------------------
# 3. Create ConfigMap and Secret from env files
# ---------------------------------------------------------------------------
echo "[deploy] Creating ConfigMap debter-config from .env.config..."
kubectl create configmap debter-config \
  --from-env-file="$REPO_ROOT/.env.config" \
  --dry-run=client -o yaml | kubectl apply -f -

echo "[deploy] Creating Secret debter-secrets from .env.dev.pg..."
kubectl create secret generic debter-secrets \
  --from-env-file="$REPO_ROOT/.env.dev.pg" \
  --dry-run=client -o yaml | kubectl apply -f -

# ---------------------------------------------------------------------------
# 4. Apply all k8s manifests
# ---------------------------------------------------------------------------
echo "[deploy] Applying k8s manifests..."
kubectl apply -f "$REPO_ROOT/k8s/"

echo "[deploy] Deployment complete."
