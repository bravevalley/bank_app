startcont:
	docker run -d --name bankpro -e POSTGRES_PASSWORD=aregbesola -e POSTGRES_USER=root -p 15432:5432 postgres

createdb:
	docker exec -it bankpro createdb --owner=root --username=root omnibank

dropdb: 
	docker exec -it bankpro dropdb omnibank

migrateup:
	migrate -path db/migration/ -database "postgres://root:aregbesola@127.0.0.1:15432/omnibank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration/ -database "postgres://root:aregbesola@127.0.0.1:15432/omnibank?sslmode=disable" -verbose down

PHONY: startcont createdb dropdb migrateup migratedown
