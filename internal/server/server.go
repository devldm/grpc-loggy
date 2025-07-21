package server

import (
	"context"
	"encoding/json"
	"fmt"
	pb "grpc-loggy/api/v1"
	"grpc-loggy/internal/storage"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type grpcServer struct {
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

func (s *grpcServer) SearchLogs(ctx context.Context, in *pb.SearchLogsRequest) (*pb.SearchLogsResponse, error) {
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

func (s *grpcServer) StreamLogs(stream pb.LoggyService_StreamLogsServer) error {
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

func (s *grpcServer) GetLogCount(ctx context.Context, in *pb.GetLogCountRequest) (*pb.GetLogCountResponse, error) {
	activeCount := len(s.activeLogs)
	totalCount := activeCount + len(s.archiveLogs)

	resp := &pb.GetLogCountResponse{
		TotalCount:  int32(totalCount),
		ActiveCount: int32(activeCount),
	}

	if in.IncludeArchive {
		archiveCount := int32(len(s.archiveLogs))
		resp.ArchiveCount = &archiveCount
	}

	return resp, nil
}

func NewGRPCServer() *grpc.Server {
	gsrv := grpc.NewServer()
	srv := grpcServer{
		redisClient: storage.NewRedisClient(),
	}
	pb.RegisterLoggyServiceServer(gsrv, &srv)
	return gsrv
}
