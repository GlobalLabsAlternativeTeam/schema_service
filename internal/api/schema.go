package api

import (
	"context"
	"fmt"
	"server/internal/domain"

	schema_service "server/proto"
)

type SchemaHandler interface {
	Create(authorID string, schemaName string, tasks []domain.Task) (domain.Schema, error)
	GetAll() ([]domain.Schema, error)
	GetByID(id string) (domain.Schema, error)
	DeleteByID(id string) error
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

func (s *SchemaServer) GetAllSchemas(ctx context.Context, req *schema_service.GetAllSchemasRequest) (*schema_service.GetAllSchemasResponse, error) {
	fmt.Println("START GetAllSchemas API")

	// Invoke SchemaHandler for fetching the schema
	schemas, err := s.SchemaHandler.GetAll()
	if err != nil {
		fmt.Println("Error calling SchemaHandler.GetAll: ", err)
		return nil, err
	}

	// Convert to gRPC objects
	var grpcSchemas []*schema_service.Schema
	for _, schema := range schemas {
		grpcSchemas = append(grpcSchemas, domain.SchemaToGRPC(&schema))
	}

	// Create and return gRPC response object
	response := &schema_service.GetAllSchemasResponse{
		Schemas: grpcSchemas,
	}

	fmt.Println("END GetAllSchemas API")
	return response, nil
}

func (s *SchemaServer) GetSchemaByID(ctx context.Context, req *schema_service.GetSchemaByIDRequest) (*schema_service.GetSchemaByIDResponse, error) {
	fmt.Println("START GetSchemaByID API")

	// Invoke SchemaHandler for fetching the schema
	schema, err := s.SchemaHandler.GetByID(req.SchemaId)
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
	err := s.SchemaHandler.DeleteByID(req.SchemaId)
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
