package api

import (
	"context"
	"fmt"
	"server/internal/domain"

	schema_service "server/proto"
)

type SchemaHandler interface {
	Create(authorID string, schemaName string, tasks []domain.Task) (domain.Schema, error)
	GetByID(id int) (domain.Schema, error)
	DeleteByID(id int) error
}

type SchemaServer struct {
	schema_service.UnimplementedSchemaServiceServer
	SchemaHandler SchemaHandler
}

func (s *SchemaServer) CreateSchema(ctx context.Context, req *schema_service.CreateSchemaRequest) (*schema_service.CreateSchemaResponse, error) {
	fmt.Println("START CreateSchema API")

	// Parse tasks from gRPC request
	var tasks []domain.Task = domain.TasksFromGRPC(req.Tasks)

	// Invoke SchemaHandler for creation
	schema, err := s.SchemaHandler.Create(req.AuthorId, req.SchemaName, tasks)
	if err != nil {
		fmt.Println("Error calling SchemaHandler.Create: ", err)
		return nil, err
	}

	// Create and return gRPC response object
	response := &schema_service.CreateSchemaResponse{
		Schema: domain.SchemaToGRPC(&schema),
	}

	fmt.Println("END CreateSchema API")
	return response, nil
}
