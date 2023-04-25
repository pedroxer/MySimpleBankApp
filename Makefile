postgres:
	docker run --name pstgr -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

sqlc:
	 /home/daniil/go/bin/sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres sqlc test