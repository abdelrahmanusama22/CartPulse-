.PHONY: run worker test vet build fmt up down sqlc
run:
	go run ./cmd/api
worker:
	go run ./cmd/worker
test:
	go test ./...
vet:
	go vet ./...
build:
	go build ./...
fmt:
	gofmt -w .
up:
	docker compose up --build -d
down:
	docker compose down
# Add the first migration and query before generation (see db/README.md).
sqlc:
	sqlc generate
