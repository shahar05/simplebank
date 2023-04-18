DB_URL=postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable
DB_URL_L=postgresql://postgres:7tghyMSFFqbZOLeF0s1m@simple-bank1.cesczdjwvygy.eu-west-1.rds.amazonaws.com:5432/simple_bank

postgres: 
	docker run --name mypostgres --network bank-net -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it mypostgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it mypostgres dropdb -U postgres simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL_L)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL_L)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL_L)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL_L)" -verbose down 1
	
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc test server
