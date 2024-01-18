#!/bin/sh

set -e

echo "Run migrations"
source /app/app.env
/app/goose/bin/goose -dir /app/migration postgres "$DB_SOURCE" up -v

echo "Start the app"
/app/main 
