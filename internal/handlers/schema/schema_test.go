package schema_test

import (
	"fmt"
	"reflect"
	"server/internal/domain"
	"server/internal/handlers/schema"
	"testing"
	"time"
)

// Mock data and variables

var now time.Time = time.Now()
var schemaId string = "schemaID"
var domainSchema domain.Schema = domain.Schema{
	SchemaID:   schemaId,
	AuthorID:   "authorID",
	SchemaName: "schemaName",
	CreatedAt:  now,
	UpdatedAt:  now,
	Tasks:      []domain.Task{},
}
var domainSchema2 domain.Schema = domain.Schema{
	SchemaID:   schemaId + "2",
	AuthorID:   "authorID2",
	SchemaName: "schemaName2",
	CreatedAt:  now,
	UpdatedAt:  now,
	Tasks:      []domain.Task{},
}

var task1 domain.Task = domain.Task{
	ID:          1,
	Level:       1,
	Name:        "Task 1",
	Status:      "NOT_STARTED",
	BlockedBy:   []int64{},
	Responsible: "Doctor1",
	TimeLimit:   3600,
	Children:    []domain.Task{},
	Comment: struct {
		Value string `json:"value"`
	}{Value: "Comment for Task 1"},
}

var task3 domain.Task = domain.Task{
	ID:          3,
	Level:       2,
	Name:        "Task 2",
	Status:      "NOT_STARTED",
	BlockedBy:   []int64{},
	Responsible: "Doctor2",
	TimeLimit:   7200,
	Children:    []domain.Task{},
	Comment: struct {
		Value string `json:"value"`
	}{Value: "Comment for Task 2"},
}

var task2 domain.Task = domain.Task{
	ID:          2,
	Level:       1,
	Name:        "Task 3",
	Status:      "NOT_STARTED",
	BlockedBy:   []int64{},
	Responsible: "Doctor3",
	TimeLimit:   5400,
	Children:    []domain.Task{task3},
	Comment: struct {
		Value string `json:"value"`
	}{Value: "Comment for Task 3"},
}

var emptyTasks []domain.Task = []domain.Task{}

// MockStorageProvider definitions

type MockStorageProvider struct{}

func (msp *MockStorageProvider) CreateSchema(authorID string, schemaName string, tasks []domain.Task) (domain.Schema, error) {
	if schemaName == "UsedSchemaName" {
		return domain.Schema{}, fmt.Errorf("schema with name '%s' already exists", schemaName)
	}

	schema := domain.Schema{
		SchemaID:   schemaId,
		AuthorID:   authorID,
		SchemaName: schemaName,
		CreatedAt:  now,
		UpdatedAt:  now,
		Tasks:      tasks,
	}

	return schema, nil
}

func (msp *MockStorageProvider) GetAllSchemas() ([]domain.Schema, error) {
	return []domain.Schema{domainSchema, domainSchema2}, nil
}

func (msp *MockStorageProvider) GetSchemaByID(id string) (domain.Schema, error) {
	if id == "NotPresentSchemaID" {
		return domain.Schema{}, fmt.Errorf("schema with id=<%s> not found", id)
	}

	schema := domain.Schema{
		SchemaID:   id,
		AuthorID:   domainSchema.AuthorID,
		SchemaName: domainSchema.SchemaName,
		CreatedAt:  now,
		UpdatedAt:  now,
		Tasks:      []domain.Task{},
	}

	return schema, nil
}

func (msp *MockStorageProvider) DeleteSchemaByID(id string) error {
	if id == "NotPresentSchemaID" {
		return fmt.Errorf("schema with id=<%s> not found", id)
	}

	return nil
}

// Tests

func TestCreate(t *testing.T) {
	mockStorageProvider := &MockStorageProvider{}
	schemaService := &schema.Schema{StorageProvider: mockStorageProvider}

	t.Run("UsedSchemaName", func(t *testing.T) {
		_, err := schemaService.Create(domainSchema.AuthorID, "UsedSchemaName", emptyTasks)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("ValidSchemaNameWithoutTasks", func(t *testing.T) {
		createdSchema, err := schemaService.Create(domainSchema.AuthorID, domainSchema.SchemaName, emptyTasks)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if createdSchema.AuthorID != domainSchema.AuthorID {
			t.Errorf("Expected AuthorID='%s', found: %s", domainSchema.AuthorID, createdSchema.AuthorID)
		}

		if createdSchema.SchemaName != domainSchema.SchemaName {
			t.Errorf("Expected SchemaName='%s', found: %s", domainSchema.SchemaName, createdSchema.SchemaName)
		}

		if len(createdSchema.Tasks) != len(emptyTasks) {
			t.Errorf("Expected %d tasks, got %d", len(emptyTasks), len(createdSchema.Tasks))
		}
	})

	t.Run("ValidSchemaNameWithTasks", func(t *testing.T) {
		expectedTasks := []domain.Task{task1, task2}
		createdSchema, err := schemaService.Create(domainSchema.AuthorID, domainSchema.SchemaName, expectedTasks)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if createdSchema.AuthorID != domainSchema.AuthorID {
			t.Errorf("Expected AuthorID='%s', found: %s", domainSchema.AuthorID, createdSchema.AuthorID)
		}

		if createdSchema.SchemaName != domainSchema.SchemaName {
			t.Errorf("Expected SchemaName='%s', found: %s", domainSchema.SchemaName, createdSchema.SchemaName)
		}

		if len(createdSchema.Tasks) != len(expectedTasks) {
			t.Errorf("Expected %d tasks, got %d", len(expectedTasks), len(createdSchema.Tasks))
		}
	})
}

func TestGetAll(t *testing.T) {
	mockStorageProvider := &MockStorageProvider{}
	schemaService := &schema.Schema{StorageProvider: mockStorageProvider}

	t.Run("GetAll", func(t *testing.T) {
		foundSchemas, err := schemaService.GetAll()
		expectedSchemas := []domain.Schema{domainSchema, domainSchema2}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(foundSchemas) != len(expectedSchemas) {
			t.Errorf("Expected len(schemas)='%d', found: %d", len(expectedSchemas), len(foundSchemas))
		}

		// Same results
		if !reflect.DeepEqual(expectedSchemas, foundSchemas) {
			t.Errorf("Expected schemas %+v, got %+v", expectedSchemas, foundSchemas)
		}

	})
}

func TestGetByID(t *testing.T) {
	mockStorageProvider := &MockStorageProvider{}
	schemaService := &schema.Schema{StorageProvider: mockStorageProvider}

	t.Run("NotPresentSchemaID", func(t *testing.T) {
		_, err := schemaService.GetByID("NotPresentSchemaID")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("PresentSchemaID", func(t *testing.T) {
		foundSchema, err := schemaService.GetByID(schemaId)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if foundSchema.SchemaID != schemaId {
			t.Errorf("Expected SchemaId='%s', found: %s", schemaId, foundSchema.SchemaID)
		}

		if foundSchema.AuthorID != domainSchema.AuthorID {
			t.Errorf("Expected AuthorID='%s', found: %s", domainSchema.AuthorID, foundSchema.AuthorID)
		}

		if foundSchema.SchemaName != domainSchema.SchemaName {
			t.Errorf("Expected SchemaName='%s', found: %s", domainSchema.SchemaName, foundSchema.SchemaName)
		}
	})
}

func TestDeleteByID(t *testing.T) {
	mockStorageProvider := &MockStorageProvider{}
	schemaService := &schema.Schema{StorageProvider: mockStorageProvider}

	t.Run("NotPresentSchemaID", func(t *testing.T) {
		err := schemaService.DeleteByID("NotPresentSchemaID")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("PresentSchemaID", func(t *testing.T) {
		err := schemaService.DeleteByID(schemaId)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
