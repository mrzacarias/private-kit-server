#!/bin/sh

set -e

cd "$(dirname "$0")/.."

echo "==> Cleaning up..."
docker-compose -f ./docker-compose.yml down --rmi=local --volumes --remove-orphans

echo "==> Running private-kit-server..."
docker-compose -f ./docker-compose.yml up --build
