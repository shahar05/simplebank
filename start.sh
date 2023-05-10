#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migration -database postgresql://postgres:7tghyMSFFqbZOLeF0s1m@simple-bank1.cesczdjwvygy.eu-west-1.rds.amazonaws.com:5432/simple_bank -verbose up


echo "start the app"
exec "$@"