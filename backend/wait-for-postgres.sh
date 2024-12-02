#!/bin/sh
# wait-for-postgres.sh

set -e

RETRIES=5

# Debug logging
echo "Checking Postgres connection: host=$DB_HOST, user=$DB_USER, db=$DB_NAME"

until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q' || [ $RETRIES -eq 0 ]; do
  echo "Waiting for postgres server to start, $((RETRIES)) remaining attempts..."
  RETRIES=$((RETRIES-=1))
  sleep 1
done

echo "Postgres is up - executing command"

# This ensures the Go application runs after Postgres is ready
exec "$@"