package main

import (
	"flag"
	"fmt"
	"grpc-loggy/internal/server"
	"log"
	"net"
)

var port = flag.Int("port", 8080, "The server port")

func main() {
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server listening on %s", addr)

	srv := server.NewGRPCServer()
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
