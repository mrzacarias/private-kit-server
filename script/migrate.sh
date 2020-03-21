#!/bin/sh

set -e

cd "$(dirname "$0")/.."

echo "==> Downloading Cloud SQL Proxy"
docker-compose run --rm web wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O cloud_sql_proxy
docker-compose run --rm web chmod +x cloud_sql_proxy

# If there is at least one parameter passed to this script
if [ "$1" != "" ]; then
  # Concatenate all parameters passed into the MIGRATE_CMD variable
  while [ "$1" != "" ]; do
    if [ "$MIGRATE_CMD" = "" ]; then
      MIGRATE_CMD="$1"
    else
      MIGRATE_CMD="$MIGRATE_CMD $1"
    fi
    shift
  done
# Otherwise, default to "up"
else
  MIGRATE_CMD="up"
fi

echo "==> Running migrations"

if [ "$PKS_ENV" = "STAGING" ]; then
  DB_NAME="private-kit-server_staging"
  INSTANCES="eq-private-kit-server-staging:us-east4:private-kit-server"
elif [ "$PKS_ENV" = "PROD" ]; then
  DB_NAME="private-kit-server_prod"
  INSTANCES="eq-private-kit-server-prod:us-east4:private-kit-server"
else
  echo "PKS_ENV must be STAGING or PROD if PKS_DB_PASS is passed"
  exit 1
fi

if [ ! -d ~/.config/gcloud ]; then
  echo "You need to have your Google Cloud credentials configured and stored at ~/.config/gcloud"
  echo "To set up your credentials, try running: gcloud auth application-default login"
  exit 1
fi

# proxy connection to Google Cloud SQL on localhost:5433
# sleep 1 to give time for the proxy to start
# run migration through the proxy
echo "==> Beginning '$MIGRATE_CMD' on '$INSTANCES' with database '$DB_NAME'"
docker-compose run -v ~/.config/gcloud:/root/.config/gcloud --rm migrate_db bash -c "echo '==> Download pq and migrate' \
  && go get -u -d github.com/lib/pq \
  && curl -OL https://github.com/golang-migrate/migrate/releases/download/v4.2.1/migrate.linux-amd64.tar.gz \
  && tar -xvzf migrate.linux-amd64.tar.gz \
  && mv migrate.linux-amd64 /usr/local/bin/migrate \
  && echo '==> Opening cloud_sql_proxy to $INSTANCES on port 5433' \
  && ./cloud_sql_proxy -instances='$INSTANCES'=tcp:5433 & sleep 10 \
  && echo '==> Running migrations' \
  && migrate -verbose -path db/migrations -database postgres://private-kit-server_migrate:'$PKS_DB_PASS'@127.0.0.1:5433/'$DB_NAME'?sslmode=disable $MIGRATE_CMD"
