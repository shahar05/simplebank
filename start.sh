#!/bin/sh

set -e

echo "run db migrations"
ls
echo "1111111"
ls /app
echo "$DB_SOURCE"
source /app/app.env
source app.env
echo $DB_SOURCE
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"