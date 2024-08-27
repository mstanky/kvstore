package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"kvstore_project/kvstore" // Adjust this to your module path

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Define command-line flags
	action := flag.String("action", "get", "Action to perform: get or set")
	key := flag.String("key", "", "Key to set or get")
	value := flag.String("value", "", "Value to set (only for set action)")
	flag.Parse()

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := kvstore.NewKeyValueStoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch *action {
	case "set":
		if *key == "" || *value == "" {
			log.Fatalf("Key and value must be provided for set action")
		}
		setResp, err := client.Set(ctx, &kvstore.SetRequest{Key: *key, Value: *value})

		if err != nil {
			log.Fatalf("Failed to set key-value: %v", err)
		}
		fmt.Printf("Set Response: %v\n", setResp.Success)

	case "get":
		if *key == "" {
			log.Fatalf("Key must be provided for get action")
		}
		getResp, err := client.Get(ctx, &kvstore.GetRequest{Key: *key})
		if err != nil {
			log.Fatalf("Failed to get value: %v", err)
		}
		if getResp.Found {
			fmt.Printf("Get Response: Key = %s, Value = %s\n", *key, getResp.Value)
		} else {
			fmt.Println("Key not found")
		}

	default:
		log.Fatalf("Unknown action: %s. Use 'get' or 'set'.", *action)
	}
}
