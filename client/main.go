package main

import (
	"context"
	"flag"
	pb "grpc-loggy/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLoggyClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SearchLogs(ctx, &pb.SearchRequest{SearchId: 10, Query: "hello"})
	if err != nil {
		log.Fatalf("Could not search logs: %v", err)
	}
	for _, cont := range r.GetLog() {
		log.Printf("Log: %s", cont.String())
	}
}
