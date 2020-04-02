#!/usr/bin/env bash
set -euo pipefail

# Requires ko & minikube
eval $(minikube docker-env)
ko publish ./cmd/eznft/ -L -B
docker tag ko.local/eznft eznft

kubectl --context=minikube delete job eznft --ignore-not-found=true
kubectl --context=minikube delete job quick --ignore-not-found=true
kubectl --context=minikube delete job slow --ignore-not-found=true
kustomize build ./resources/minikube/ | kubectl --context=minikube apply -f -