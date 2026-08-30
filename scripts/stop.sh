#!/usr/bin/env bash
# =============================================================================
# Debter — Stop or delete the local Minikube cluster
# =============================================================================
# Subcommands:
#   stop   — halt the Minikube VM (preserves cluster + PVC data)
#   delete — destroy the Minikube cluster (removes pods, PVC data, etc.)
#
# Usage: ./scripts/stop.sh stop
#        ./scripts/stop.sh delete
# =============================================================================

set -euo pipefail

command -v minikube >/dev/null 2>&1 \
  || { echo "[stop] ERROR: minikube is not installed or not on PATH. Exiting." >&2; exit 1; }

case "${1:-}" in
  stop)
    echo "[stop] Stopping minikube..."
    minikube stop
    ;;
  delete)
    echo "[stop] Deleting minikube cluster..."
    minikube delete --all
    ;;
  *)
    echo "Usage: $0 {stop|delete}" >&2
    exit 1
    ;;
esac

echo "[stop] Done."
