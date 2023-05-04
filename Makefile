postgres:
	docker run --name pstgr -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

sqlc:
	 /home/daniil/go/bin/sqlc generate
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/melbank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/melbank?sslmode=disable" -verbose down
test:
	go test -v -cover ./...
server:
	go run main.go

.PHONY: postgres sqlc test migrateup migratedown server
