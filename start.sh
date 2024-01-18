#!/bin/sh

set -e

set -a
. /app/app.env 
set +a

echo "Run migrations"
/app/goose/bin/goose -dir /app/migration postgres "$DB_SOURCE" up -v

echo "Start the app"
/app/main 
