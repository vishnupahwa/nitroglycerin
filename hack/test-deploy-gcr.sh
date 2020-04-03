#!/usr/bin/env bash
set -euo pipefail

export KO_DOCKER_REPO=gcr.io/phoenix-sandbox-one

CONTEXT=$1
# Requires ko & minikube
ko publish ./cmd/eznft/ -B

kubectl --context=${CONTEXT} delete job eznft --ignore-not-found=true
kubectl --context=${CONTEXT} delete job quick --ignore-not-found=true
kubectl --context=${CONTEXT} delete job slow --ignore-not-found=true
kustomize build ./resources/dev | kubectl --context=${CONTEXT} apply -f -