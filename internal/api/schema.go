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

func (s *SchemaServer) GetSchemaByID(ctx context.Context, req *schema_service.GetSchemaByIDRequest) (*schema_service.GetSchemaByIDResponse, error) {
	fmt.Println("START GetSchemaByID API")

	// Invoke SchemaHandler for fetching the schema
	schema, err := s.SchemaHandler.GetByID(int(req.SchemaId))
	if err != nil {
		fmt.Println("Error calling SchemaHandler.GetByID: ", err)
		return nil, err
	}

	// Create and return gRPC response object
	response := &schema_service.GetSchemaByIDResponse{
		SchemaId: req.SchemaId,
		Schema:   domain.SchemaToGRPC(&schema),
	}

	fmt.Println("END GetSchemaByID API")
	return response, nil
}

func (s *SchemaServer) DeleteSchemaByID(ctx context.Context, req *schema_service.DeleteSchemaByIDRequest) (*schema_service.DeleteSchemaByIDResponse, error) {
	fmt.Println("START DeleteSchemaByID API")

	// Invoke SchemaHandler for deleting the schema
	err := s.SchemaHandler.DeleteByID(int(req.SchemaId))
	if err != nil {
		fmt.Println("Error calling SchemaHandler.DeleteByID: ", err)
		return nil, err
	}

	// Create and return gRPC response object
	response := &schema_service.DeleteSchemaByIDResponse{
		SchemaId: req.SchemaId,
	}

	fmt.Println("END DeleteSchemaByID API")
	return response, nil
}
