postgres: 
	docker run --name mypostgres -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it mypostgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it mypostgres dropdb -U postgres simple_bank

migrateup:
	migrate -path db/migration -database postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable -verbose up

migratedown:
	migrate -path db/migration -database postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc
