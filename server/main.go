package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"grpc-loggy/redispkg"
	"log"
	"net"
	"strings"

	pb "grpc-loggy/proto/v1"

	"github.com/redis/go-redis/v9"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

type server struct {
	pb.UnimplementedLoggyServiceServer
	activeLogs  []*pb.Log
	archiveLogs []*pb.Log
	redisClient *redis.Client
}

func findSubstring(sliceStrings []*pb.Log, substring string) []*pb.Log {
	substring = strings.ToLower(substring)
	var matches []*pb.Log
	for _, v := range sliceStrings {
		if strings.Contains(strings.ToLower(v.Content), substring) {
			matches = append(matches, v)
		}
	}
	return matches
}

func (s *server) SearchLogs(ctx context.Context, in *pb.SearchLogsRequest) (*pb.SearchLogsResponse, error) {
	log.Printf("Received %v", in.GetQuery())
	cachedState, err := s.redisClient.HGet(ctx, string(in.GetQuery()), "logs").Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			log.Printf("HGet failed on cache key: %s, error: %v", in.GetQuery(), err)
		}
	}

	if len(cachedState) > 0 {
		fmt.Println("returning cached values")
		var cachedData []*pb.Log
		err := json.Unmarshal([]byte(cachedState), &cachedData)
		if err != nil {
			log.Printf("Failed to parse cached json: %v", err)
		}
		return &pb.SearchLogsResponse{TotalCount: int32(len(cachedData)), Log: cachedData}, nil
	}
	matches := findSubstring(s.activeLogs, in.GetQuery())
	var res []*pb.Log
	for _, l := range matches {
		res = append(res, &pb.Log{Content: l.Content, Level: l.Level, Origin: l.Origin})
	}

	cacheMe, err := json.Marshal(res)
	if err != nil {
		log.Printf("Failed to marshal api response: %v", err)
	}

	err = s.redisClient.HSet(ctx, in.GetQuery(), "logs", cacheMe).Err()
	if err != nil {
		log.Printf("HSet failed: %v", err)
	}
	return &pb.SearchLogsResponse{TotalCount: int32(len(res)), Log: res}, nil
}

func (s *server) StreamLogs(stream pb.LoggyService_StreamLogsServer) error {
	for {
		msg, err := stream.Recv()
		// TODO: handle EOF here first
		if err != nil {
			return err
		}

		if len(s.activeLogs) != 10000 {
			s.activeLogs = append(s.activeLogs, msg.Log)
		} else {
			s.archiveLogs = s.activeLogs
			s.activeLogs = []*pb.Log{}
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLoggyServiceServer(s, &server{redisClient: redispkg.Connect()})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
