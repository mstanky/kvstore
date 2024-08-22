package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"kvstore_project/kvstore" // Adjust this to your module path
)

type kvServer struct {
	kvstore.UnimplementedKeyValueStoreServer
	store map[string]string
}

func (s *kvServer) Get(ctx context.Context, req *kvstore.GetRequest) (*kvstore.GetResponse, error) {
	value, found := s.store[req.Key]
	return &kvstore.GetResponse{
		Value: value,
		Found: found,
	}, nil
}

func (s *kvServer) Set(ctx context.Context, req *kvstore.SetRequest) (*kvstore.SetResponse, error) {
	s.store[req.Key] = req.Value
	return &kvstore.SetResponse{
		Success: true,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	kvstore.RegisterKeyValueStoreServer(s, &kvServer{store: make(map[string]string)})

	fmt.Println("gRPC server is running on port 50051...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
