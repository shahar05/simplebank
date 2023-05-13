#!/bin/sh

set -e

echo "run db migrations"
source /app/app.env
echo "before migrate"
cat /app/app.env
ls /app
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up
echo "after migrate"

echo "start the app"
exec "$@"