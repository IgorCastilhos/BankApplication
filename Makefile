DB_URL=postgresql://root:secret@localhost:5432/bank?sslmode=disable

postgres:
	docker run --name postgres16.1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.1-alpine

docker_start:
	docker start postgres16.1

docker_stop:
	docker stop postgres16.1

createdb:
	docker exec -it postgres16.1 createdb --username=root --owner=root bank

dropdb:
	docker exec -it postgres16.1 dropdb bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go  github.com/IgorCastilhos/BankApplication/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown test docker_start docker_stop server mock