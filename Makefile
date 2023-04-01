runcont:
	docker run -d --name bankpro -e POSTGRES_PASSWORD=aregbesola -e POSTGRES_USER=root -p 15432:5432 postgres

startcont:
	docker start bankpro

createdb:
	docker exec -it bankpro createdb --owner=root --username=root omnibank

dropdb: 
	docker exec -it bankpro dropdb omnibank

migrateup:
	migrate -path db/migration/ -database "postgres://root:aregbesola@127.0.0.1:15432/omnibank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration/ -database "postgres://root:aregbesola@127.0.0.1:15432/omnibank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run .

PHONY: startcont createdb dropdb migrateup migratedown sqlc test runcont server
