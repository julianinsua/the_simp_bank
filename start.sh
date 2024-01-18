#!/bin/sh

set -e

echo "Run migrations"
set -o allexport
source /app/app.env 
set +o allexport
/app/goose/bin/goose -dir /app/migration postgres "$DB_SOURCE" up -v

echo "Start the app"
/app/main 
