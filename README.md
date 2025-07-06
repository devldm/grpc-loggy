# grpc-loggy

A lightweight gRPC-based logging server and client system written in Go.

## Overview

grpc-loggy enables applications to stream, store, and search logs over gRPC. It provides:

- A gRPC server that accepts log streams, stores them, and allows searching via gRPC methods.
- A simple client for sending and searching logs using the defined gRPC interface.

## Features

- **Stream Logs:** Clients can stream log messages to the server.
- **Search Logs:** Search log entries by substring query.
- **Automatic Log Archival:** Moves active logs to an archive after a threshold.

## Directory Structure

- `server/` &mdash; gRPC server implementation.
- `client/` &mdash; Example gRPC client usage.
- `proto/` &mdash; Protobuf service definitions and message types.

## Usage

### Server

```sh
cd server
go run main.go -port=50051
```

- The server listens on the specified TCP port (default: 50051).
- Accepts log entries via `StreamLogs`.
- Search logs using `SearchLogs`.

### Client

```sh
cd client
go run main.go -addr=localhost:50051
```

- The client connects to the server and issues a sample search request.
- Modify the client to stream logs or perform other actions as needed.

## API

Loggy uses gRPC with the following main RPCs:

- `StreamLogs` &mdash; Client-side streaming for log ingestion.
- `SearchLogs` &mdash; Unary call for searching logs by substring.

## Example

- Start the server.
- Run the client. By default, the client sends a search request with query "hello".

## Requirements

- Go 1.18+
- [gRPC Go](https://github.com/grpc/grpc-go)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)

## Customization

- Update the proto file in `proto/` to add fields or RPCs.
- Extend server/client logic as needed for your environment.
