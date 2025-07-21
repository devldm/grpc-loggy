# grpc-loggy

A lightweight gRPC-based logging server and client system written in Go.

## Overview

grpc-loggy enables applications to stream, store, and search logs over gRPC. It provides:

- A gRPC server that accepts log streams, stores them in Redis, and allows querying
- A simple client for sending and searching logs

## Features

- **Stream Logs:** Client-side streaming for log ingestion
- **Get Log Count:** Track active, archived, and total logs
- **Search Logs**: Search log entries by substring query.
- **Redis Storage:** Fast log persistence
- **Automatic Archival:** Moves logs to archive after 10,000 entries

## Usage

```bash
# Start server (default port 8080)
just server

# Run client
just client

# See all commands
just
```

**Manual commands:**

```bash
go run ./cmd/server
go run ./cmd/client
```

## API

Main gRPC methods:

- `StreamLogs` — Client-side streaming for log ingestion
- `GetLogCount` — Get counts of active and archived logs
- `SearchLogs` - Search the stored logs

## Structure

```
├── cmd/           # Server and client applications
├── internal/      # Server logic & Redis storage
├── api/v1/        # Protobuf definitions and generated code
└── justfile       # Build and run commands
```

## Requirements

- Go 1.18+
- Redis server
- protoc (for development)

## Development

```bash
just proto    # Generate protobuf files
just build    # Build binaries
```
