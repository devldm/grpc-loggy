package main

import (
	"context"
	"flag"
	pb "grpc-loggy/proto/v1"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func seedLogs() []*pb.Log {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	level := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	service := []string{"grpc-loggy", "data-svc", "db-svc", "err-handler", "s3-archive-svc"}
	content := []string{"Successfully processed", "Succeeded with 1 warning", "Failed to get x", "Failure in x svc", "Job took x seconds"}
	var logs []*pb.Log

	for range 10000 {
		randLog := &pb.Log{Content: content[r.Intn((len(content)))], Level: level[r.Intn(len(level))], Origin: service[r.Intn(len(service))]}
		logs = append(logs, randLog)
	}

	return logs
}

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLoggyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	logStream := seedLogs()
	stream, err := c.StreamLogs(ctx)
	if err != nil {
		log.Fatalf("failed to set up stream")
	}
	for _, logEntry := range logStream {
		stream.Send(&pb.StreamLogsRequest{
			Log: &pb.Log{
				Content: logEntry.Content,
				Level:   logEntry.Level,
				Origin:  logEntry.Origin,
			},
		})
	}

	stream.CloseAndRecv()
	r, err := c.SearchLogs(ctx, &pb.SearchLogsRequest{SearchId: 10, Query: "failed"})
	if err != nil {
		log.Fatalf("Could not search logs: %v", err)
	}
	for _, cont := range r.GetLog() {
		log.Printf("Log: %s", cont.String())
	}
}
