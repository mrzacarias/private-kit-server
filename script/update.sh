#!/bin/sh

set -e

cd "$(dirname "$0")/.."

echo "==> Ensuring packages..."
docker-compose run --rm web dep ensure -v
