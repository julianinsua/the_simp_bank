#!/bin/sh

set -e

set -o allexport
source /app/app.env 
set +o allexport

echo "Run migrations"
/app/goose/bin/goose -dir /app/migration postgres "$DB_SOURCE" up -v

echo "Start the app"
/app/main 
