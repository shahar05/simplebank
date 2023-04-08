DB_URL=postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable


postgres: 
	docker run --name mypostgres -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it mypostgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it mypostgres dropdb -U postgres simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
	
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc test server
