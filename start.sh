#!/bin/sh

set -e

if [ -f "/app/app.env" ]; then
  # Source and export variables
  source "app.env"
  export $(cut -d= -f1 /app/app.env)
  echo "Variables from app.env have been sourced and exported."
else
  echo "Error: app.env file not found."
fi

echo "Run migrations"
/app/goose/bin/goose -dir /app/migration postgres "$DB_SOURCE" up -v

echo "Start the app"
/app/main 
