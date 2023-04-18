#!/bin/sh

set -e

echo "run db migrations"
ls
ls /app
echo "asdasd"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"