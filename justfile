# Go project justfile

# Show available commands
default:
    @just --list

# Run the server
server:
    go run ./cmd/server

# Run the client
client:
    go run ./cmd/client

# Build project - output to bin/
build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client

# Compile proto files
proto:
	protoc --go_out=./api/v1/ --go-grpc_out=./api/v1 api/v1/*.proto
