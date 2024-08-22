package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"kvstore_project/kvstore" // Adjust this to your module path

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := kvstore.NewKeyValueStoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Set a key-value pair
	setResp, err := client.Set(ctx, &kvstore.SetRequest{Key: "example", Value: "Hello, gRPC!"})
	if err != nil {
		log.Fatalf("Failed to set key-value: %v", err)
	}
	fmt.Printf("Set Response: %v\n", setResp.Success)

	// Get the value for the key
	getResp, err := client.Get(ctx, &kvstore.GetRequest{Key: "example"})
	if err != nil {
		log.Fatalf("Failed to get value: %v", err)
	}
	if getResp.Found {
		fmt.Printf("Get Response: Key = example, Value = %s\n", getResp.Value)
	} else {
		fmt.Println("Key not found")
	}
}
