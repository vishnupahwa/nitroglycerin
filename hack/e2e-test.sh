#!/usr/bin/env bash
set -euo pipefail

# Requires ko & kind
# Run setup-e2e.sh before this
kubectl --context=kind-eznft delete jobs --all
kubectl --context=kind-eznft delete pods --all

go test ./internal/test/... --tags=e2e -count=1