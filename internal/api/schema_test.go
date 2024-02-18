package api_test

import (
	"context"
	"fmt"
	"server/internal/api"
	"server/internal/domain"
	schema_service "server/proto"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Mock data and variables

var now time.Time = time.Now()
var schema_id string = "schemaID"
var domain_schema domain.Schema = domain.Schema{
	SchemaID:   schema_id,
	AuthorID:   "authorID",
	SchemaName: "schemaName",
	CreatedAt:  now,
	UpdatedAt:  now,
	Tasks:      []domain.Task{},
}
var grpc_schema schema_service.Schema = *domain.SchemaToGRPC(&domain_schema)

var task1 schema_service.Task = schema_service.Task{
	Id:          1,
	Level:       1,
	Name:        "Task 1",
	Status:      schema_service.TaskStatus_TASK_STATUS_NOT_STARTED,
	BlockedBy:   []int64{},
	Responsible: "Doctor1",
	TimeLimit:   3600,
	Children:    []*schema_service.Task{},
	Comment:     wrapperspb.String("Task 1 comment"),
}

var task3 schema_service.Task = schema_service.Task{
	Id:          3,
	Level:       2,
	Name:        "Task 3",
	Status:      schema_service.TaskStatus_TASK_STATUS_NOT_STARTED,
	BlockedBy:   []int64{},
	Responsible: "Doctor2",
	TimeLimit:   3600,
	Children:    []*schema_service.Task{},
	Comment:     wrapperspb.String("Task 3 comment"),
}

var task2 schema_service.Task = schema_service.Task{
	Id:          2,
	Level:       1,
	Name:        "Task 2",
	Status:      schema_service.TaskStatus_TASK_STATUS_NOT_STARTED,
	BlockedBy:   []int64{},
	Responsible: "Doctor2",
	TimeLimit:   3600,
	Children:    []*schema_service.Task{&task3},
	Comment:     wrapperspb.String("Task 2 comment"),
}

// MockSchemaHandler definitions

type MockSchemaHandler struct{}

func (msh *MockSchemaHandler) Create(authorID string, schemaName string, tasks []domain.Task) (domain.Schema, error) {
	if schemaName == "UsedSchemaName" {
		return domain.Schema{}, fmt.Errorf("schema with name '%s' already exists", schemaName)
	}

	schema := domain.Schema{
		SchemaID:   "schemaID",
		AuthorID:   authorID,
		SchemaName: schemaName,
		CreatedAt:  now,
		UpdatedAt:  now,
		Tasks:      tasks,
	}

	return schema, nil
}

func (msh *MockSchemaHandler) GetByID(id string) (domain.Schema, error) {
	if id == "NotPresentSchemaID" {
		return domain.Schema{}, fmt.Errorf("schema with id=<%s> not found", id)
	}

	schema := domain.Schema{
		SchemaID:   id,
		AuthorID:   domain_schema.AuthorID,
		SchemaName: domain_schema.SchemaName,
		CreatedAt:  now,
		UpdatedAt:  now,
		Tasks:      []domain.Task{},
	}

	return schema, nil
}

func (msh *MockSchemaHandler) DeleteByID(id string) error {
	if id == "NotPresentSchemaID" {
		return fmt.Errorf("schema with id=<%s> not found", id)
	}

	return nil
}

// Tests

func TestCreateSchema(t *testing.T) {
	// Create instances of your dependencies (handlers, storage, etc.)
	mockHandler := &MockSchemaHandler{}
	apiHandler := api.SchemaServer{SchemaHandler: mockHandler}

	t.Run("UsedSchemaName", func(t *testing.T) {
		request := schema_service.CreateSchemaRequest{
			AuthorId:   "authorID",
			SchemaName: "UsedSchemaName",
			Tasks:      []*schema_service.Task{},
		}
		_, err := apiHandler.CreateSchema(context.Background(), &request)

		if err == nil {
			t.Errorf("Expected error, got nil")
			return
		}
	})

	t.Run("ValidSchemaNameWithoutTasks", func(t *testing.T) {
		expectedAuthorId := "authorID"
		expectedSchemaName := "ValidSchemaName"
		expectedTasks := []*schema_service.Task{}

		request := schema_service.CreateSchemaRequest{
			AuthorId:   expectedAuthorId,
			SchemaName: expectedSchemaName,
			Tasks:      expectedTasks,
		}
		response, err := apiHandler.CreateSchema(context.Background(), &request)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		schema := response.Schema

		if schema.AuthorId != expectedAuthorId {
			t.Errorf("Expected AuthorID='%s', found: %s", expectedAuthorId, schema.AuthorId)
		}

		if schema.SchemaName != expectedSchemaName {
			t.Errorf("Expected SchemaName='%s', found: %s", expectedSchemaName, schema.SchemaName)
		}

		if len(expectedTasks) != len(schema.Tasks) {
			t.Errorf("Expected %d tasks, got %d", len(expectedTasks), len(schema.Tasks))
		}
	})

	t.Run("ValidSchemaNameWithTasks", func(t *testing.T) {
		expectedAuthorId := "authorID"
		expectedSchemaName := "ValidSchemaName"
		expectedTasks := []*schema_service.Task{&task1, &task2}

		request := schema_service.CreateSchemaRequest{
			AuthorId:   expectedAuthorId,
			SchemaName: expectedSchemaName,
			Tasks:      expectedTasks,
		}
		response, err := apiHandler.CreateSchema(context.Background(), &request)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		schema := response.Schema

		if schema.AuthorId != expectedAuthorId {
			t.Errorf("Expected AuthorID='%s', found: %s", expectedAuthorId, schema.AuthorId)
		}

		if schema.SchemaName != expectedSchemaName {
			t.Errorf("Expected SchemaName='%s', found: %s", expectedSchemaName, schema.SchemaName)
		}

		if len(expectedTasks) != len(schema.Tasks) {
			t.Errorf("Expected %d tasks, got %d", len(expectedTasks), len(schema.Tasks))
		}
	})
}

func TestGetSchemaByID(t *testing.T) {
	// Create instances of your dependencies (handlers, storage, etc.)
	mockHandler := &MockSchemaHandler{}
	apiHandler := api.SchemaServer{SchemaHandler: mockHandler}

	t.Run("NotPresentSchemaID", func(t *testing.T) {
		request := schema_service.GetSchemaByIDRequest{
			SchemaId: "NotPresentSchemaID",
		}
		_, err := apiHandler.GetSchemaByID(context.Background(), &request)

		if err == nil {
			t.Errorf("Expected error, got nil")
			return
		}
	})

	t.Run("PresentSchemaID", func(t *testing.T) {
		request := schema_service.GetSchemaByIDRequest{
			SchemaId: schema_id,
		}
		response, err := apiHandler.GetSchemaByID(context.Background(), &request)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if schema_id != response.SchemaId {
			t.Errorf("Expected SchemaId='%s', found: %s", schema_id, response.SchemaId)
		}

		schema := response.Schema

		if schema.AuthorId != grpc_schema.AuthorId {
			t.Errorf("Expected AuthorID='%s', found: %s", grpc_schema.AuthorId, schema.AuthorId)
		}

		if schema.SchemaName != grpc_schema.SchemaName {
			t.Errorf("Expected SchemaName='%s', found: %s", grpc_schema.SchemaName, schema.SchemaName)
		}
	})
}

func TestDeleteSchemaByID(t *testing.T) {
	// Create instances of your dependencies (handlers, storage, etc.)
	mockHandler := &MockSchemaHandler{}
	apiHandler := api.SchemaServer{SchemaHandler: mockHandler}

	t.Run("NotPresentSchemaID", func(t *testing.T) {
		request := schema_service.DeleteSchemaByIDRequest{
			SchemaId: "NotPresentSchemaID",
		}
		_, err := apiHandler.DeleteSchemaByID(context.Background(), &request)

		if err == nil {
			t.Errorf("Expected error, got nil")
			return
		}
	})

	t.Run("PresentSchemaID", func(t *testing.T) {
		request := schema_service.DeleteSchemaByIDRequest{
			SchemaId: schema_id,
		}
		response, err := apiHandler.DeleteSchemaByID(context.Background(), &request)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if schema_id != response.SchemaId {
			t.Errorf("Expected SchemaId='%s', found: %s", schema_id, response.SchemaId)
		}
	})
}
