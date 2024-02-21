package storage_test

import (
	"reflect"
	"server/internal/domain"
	"server/internal/providers/storage"
	"testing"
)

func TestCreateSchema(t *testing.T) {
	storageService, err := storage.NewStorage("./test_storage.json", true)
	if err != nil {
		t.Errorf("Failed to create storage: %v", err)
	}

	t.Run("Cannot create schema with used name", func(t *testing.T) {
		_, err := storageService.CreateSchema("authorID", "Schema2", []domain.Task{})

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("Creates schema with valid name", func(t *testing.T) {
		// Create schema
		createdSchema, err := storageService.CreateSchema("authorID1", "schemaName1", []domain.Task{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Schema exists after creation
		foundSchema, err := storageService.GetSchemaByID(createdSchema.SchemaID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Both are the same schemas
		if !reflect.DeepEqual(createdSchema, foundSchema) {
			t.Errorf("Expected schema %+v, got %+v", createdSchema, foundSchema)
		}
	})
}

func TestGetAllSchemas(t *testing.T) {
	storageService, err := storage.NewStorage("./test_storage.json", true)
	if err != nil {
		t.Errorf("Failed to create storage: %v", err)
	}

	t.Run("All schemas", func(t *testing.T) {
		expectedSchemasLen := 3
		foundSchemas, err := storageService.GetAllSchemas()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(foundSchemas) != expectedSchemasLen {
			t.Errorf("Expected len(schemas)='%d', found: %d", expectedSchemasLen, len(foundSchemas))
		}
	})
}

func TestGetSchemaByID(t *testing.T) {
	storageService, err := storage.NewStorage("./test_storage.json", true)
	if err != nil {
		t.Errorf("Failed to create storage: %v", err)
	}

	t.Run("Schema not present", func(t *testing.T) {
		_, err := storageService.GetSchemaByID("SchemaNotPresent")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("Schema present", func(t *testing.T) {
		schemaId := "dd19b4f7-a4be-4ec8-a48b-be6c4f769b4b"
		expectedSchemaAuthorId := "Author1"
		expectedSchemaName := "Schema2"
		foundSchema, err := storageService.GetSchemaByID(schemaId)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if foundSchema.SchemaID != schemaId {
			t.Errorf("Expected SchemaId='%s', found: %s", schemaId, foundSchema.SchemaID)
		}

		if foundSchema.AuthorID != expectedSchemaAuthorId {
			t.Errorf("Expected AuthorID='%s', found: %s", expectedSchemaAuthorId, foundSchema.AuthorID)
		}

		if foundSchema.SchemaName != expectedSchemaName {
			t.Errorf("Expected SchemaName='%s', found: %s", expectedSchemaName, foundSchema.SchemaName)
		}
	})
}

func TestDeleteSchemaByID(t *testing.T) {
	storageService, err := storage.NewStorage("./test_storage.json", true)
	if err != nil {
		t.Errorf("Failed to create storage: %v", err)
	}

	t.Run("Schema not present", func(t *testing.T) {
		err := storageService.DeleteSchemaByID("SchemaNotPresent")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("Schema present", func(t *testing.T) {
		schemaId := "dd19b4f7-a4be-4ec8-a48b-be6c4f769b4b"
		err := storageService.DeleteSchemaByID(schemaId)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Schema should not be in the storage anymore
		_, err = storageService.GetSchemaByID(schemaId)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
