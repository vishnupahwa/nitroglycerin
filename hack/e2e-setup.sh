#!/usr/bin/env bash
set -euo pipefail

kind delete cluster --name eznft --quiet
kind create cluster --name eznft
ko publish ./internal/test/e2e/testdata/hits -L -B
docker tag ko.local/hits hits
kind load --name eznft docker-image hits
kubectl apply -f ./internal/test/e2e/testdata/hits/app.yaml

ko publish ./cmd/eznft/ -L -B
docker tag ko.local/eznft eznft
kind load --name eznft docker-image eznft