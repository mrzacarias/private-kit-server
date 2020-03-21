#!/bin/sh

set -e

cd "$(dirname "$0")/.."

[ -z "$DEBUG" ] || set -x

if [ -n "$1" ]; then
  testdir="/$*"
fi

echo "==> Cleaning up..."
docker-compose -f ./docker-compose-test.yml down --rmi=local --volumes --remove-orphans

echo "==> Running private-kit-server tests..."
docker-compose -f ./docker-compose-test.yml up --build -d
docker-compose -f ./docker-compose-test.yml run --rm test go test -p 1 -cover .$testdir/...
