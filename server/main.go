package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	pb "grpc-loggy/proto"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

type server struct {
	pb.UnimplementedLoggyServer
	activeLogs  []*pb.Log
	archiveLogs []*pb.Log
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

func (s *server) SearchLogs(_ context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("Received %v", in.GetQuery())
	matches := findSubstring(s.activeLogs, in.GetQuery())
	var res []*pb.Log
	for _, l := range matches {
		res = append(res, &pb.Log{Content: l.Content, Level: l.Level, Origin: l.Origin})
	}
	return &pb.SearchResponse{TotalCount: int32(len(res)), Log: res}, nil
}

func (s *server) StreamLogs(stream pb.Loggy_StreamLogsServer) error {
	for {
		msg, err := stream.Recv()
		// TODO: handle EOF here first
		if err != nil {
			return err
		}

		if len(s.activeLogs) == 10000 {
			s.archiveLogs = s.activeLogs
			s.activeLogs = []*pb.Log{}
		} else {
			s.activeLogs = append(s.activeLogs, msg)
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
	pb.RegisterLoggyServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
