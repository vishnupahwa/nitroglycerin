#!/usr/bin/env bash
set -euo pipefail

# Requires ko & kind
# kind create cluster
ko publish ./cmd/eznft/ -L -B
docker tag ko.local/eznft eznft
kind load docker-image eznft

kubectl --context=kind-kind delete job eznft --ignore-not-found=true
kubectl --context=kind-kind delete job quick --ignore-not-found=true
kubectl --context=kind-kind delete job slow --ignore-not-found=true
kustomize build ./resources/kind/ | kubectl --context=kind-kind apply -f -