#!/usr/bin/env bash
set -euo pipefail

CONTEXT=$1

SHORT_SHA=$(git rev-parse --short HEAD)
export KO_DOCKER_REPO=registry.tools.cosmic.sky/phoenix/test

# Requires ko, kustomize & kubectl
ko publish ./cmd/eznft/ -B -t ${SHORT_SHA} -t latest

kubectl --context=${CONTEXT} delete job eznft --ignore-not-found=true
kubectl --context=${CONTEXT} delete job quick --ignore-not-found=true
kubectl --context=${CONTEXT} delete job slow --ignore-not-found=true

cd ./resources/dev

kustomize edit set image eznft=registry.tools.cosmic.sky/phoenix/test/eznft:${SHORT_SHA}
kustomize build | kubectl --context=${CONTEXT} apply -f -