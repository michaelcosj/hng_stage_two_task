entry = cmd/api/main.go

build:
	@echo "building application"
	@sqlc generate
	@go build -o bin/main $(entry)

run:
	@echo "running application"
	@go run $(entry)

test:
	@echo "testing application"
	@go test ./test -v

clean:
	@echo "cleaning binary"
	@rm -f bin/main

watch:
	@echo "watching"
	@air

migrate:
	@echo "running migrations"
	@tern migrate -m ./db/migrations

.PHONY: build run test clean
