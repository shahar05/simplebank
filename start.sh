#!/bin/sh

set -e

echo "run db migrations"
ls
echo "1111111"
ls /app
echo "$DB_SOURCE"
source /app/app.env
source app.env
echo "222222"
cat app.env
echo $DB_SOURCE
/app/migrate -path /app/migration -database postgresql://postgres:7tghyMSFFqbZOLeF0s1m@simple-bank1.cesczdjwvygy.eu-west-1.rds.amazonaws.com:5432/simple_bank -verbose up

echo "start the app"
exec "$@"