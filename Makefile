.PHONY: run build test fmt vet

# Load .env (if present) and start the server
run:
	@set -a; [ -f .env ] && . ./.env; set +a; go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
