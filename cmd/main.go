// cmd/main.go

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// Create a listener on TCP port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	server := grpc.NewServer()

	// Serve and listen for incoming requests
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
