// cmd/main.go

package main

import (
	"log"
	"net"
	"server/internal/api"
	"server/internal/handlers/schema"
	"server/internal/providers/storage"
	schema_service "server/proto"

	"google.golang.org/grpc"
)

func main() {
	// Create a listener on TCP port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create instances of your dependencies (handlers, storage, etc.)
	storageService, err := storage.NewStorage("./data/storage.json")
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	schemaHandler := &schema.Schema{StorageProvider: storageService}
	apiService := &api.SchemaServer{SchemaHandler: schemaHandler}

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the ProcessExecutionService server
	schema_service.RegisterSchemaServiceServer(server, apiService)

	// Serve and listen for incoming requests
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
