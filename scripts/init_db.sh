#!/usr/bin/env bash

set -x
set -eo pipefail

script_dir=$(pwd)

cd "$(dirname "$0")"
docker compose -f ../docker-compose.yml up -d

cd "$script_dir"

sleep 1

if ! [ -x "$(command -v migrate)" ]; then
    echo >&2 "Error: migrate is not installed."
    echo >&2 "Use:"
    echo >&2 "go get -u github.com/golang-migrate/migrate/v4/cmd/migrate"
    echo >&2 "to install it."
    exit 1
fi

DB_USER="${POSTGRES_USER:=user}"
DB_PASSWORD="${POSTGRES_PASSWORD:=password}"
DB_NAME="${POSTGRES_DB:=indegodb}"
DB_PORT="${POSTGRES_PORT:=5432}"
DB_HOST="${POSTGRES_HOST:=localhost}"

export PGPASSWORD="${DB_PASSWORD}"
until psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "postgres" -c '\q'; do
    >&2 echo "Postgres is still unavailable - sleeping"
    sleep 1
done

>&2 echo "Postgres is up and running on port ${DB_PORT} - running migration now!"

DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}
export DATABASE_URL
sqlx database create
sqlx migrate run

>&2 echo "Postgres has been migrated, ready to go!"