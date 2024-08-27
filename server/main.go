package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"kvstore_project/kvstore" // Adjust this to your module path

	"google.golang.org/grpc"
)

const dataFile = "kvstore_data.json"

type kvServer struct {
	kvstore.UnimplementedKeyValueStoreServer
	store *kvstore.KVStore
}

func (s *kvServer) Get(ctx context.Context, req *kvstore.GetRequest) (*kvstore.GetResponse, error) {
	value, found := s.store.Get(req.Key)
	return &kvstore.GetResponse{
		Value: value,
		Found: found,
	}, nil
}

func (s *kvServer) Set(ctx context.Context, req *kvstore.SetRequest) (*kvstore.SetResponse, error) {
	s.store.Set(req.Key, req.Value)
	log.Printf("Server: Set key=%s, value=%s\n", req.Key, req.Value)
	return &kvstore.SetResponse{
		Success: true,
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	store := kvstore.NewKVStore()

	//Load existing data from file
	if err := store.LoadFromFile(dataFile); err != nil {
		log.Fatalf("Failed to load data from file: %v", err)
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	kvstore.RegisterKeyValueStoreServer(s, &kvServer{store: store})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Saving data to file before shutting down...")
		if err := store.SaveToFile(dataFile); err != nil {
			log.Fatalf("Failed to save data to file: %v", err)
		}
		os.Exit(0)
	}()

	log.Println("gRPC server is running on port 50051...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
