package schema

import (
	"fmt"
	"server/internal/domain"
)

type StorageInterface interface {
	CreateSchema(authorID string, schemaName string, tasks []domain.Task) (domain.Schema, error)
	GetSchemaByID(id int) (domain.Schema, error)
	DeleteSchemaByID(id int) error
}

type Schema struct {
	StorageProvider StorageInterface
}

func (s *Schema) Create(authorID string, schemaName string, tasks []domain.Task) (domain.Schema, error) {
	fmt.Println("START Schema.Create handler")

	// Forward creation to Storage
	schema, err := s.StorageProvider.CreateSchema(authorID, schemaName, tasks)
	if err != nil {
		fmt.Println("Error creating Schema: ", err)
	}

	fmt.Println("END Schema.Create handler")
	return schema, err
}

func (s *Schema) GetByID(id int) (domain.Schema, error) {
	fmt.Println("START Schema.GetByID handler")

	// Forward fetch to Storage
	schema, err := s.StorageProvider.GetSchemaByID(id)
	if err != nil {
		fmt.Printf("Error getting Schema with id=<%d>: %s\n", id, err)
	}

	fmt.Println("END Schema.GetByID handler")
	return schema, err
}

func (s *Schema) DeleteByID(id int) error {
	fmt.Println("START Schema.DeleteByID handler")

	// Forward deletion to Storage
	err := s.StorageProvider.DeleteSchemaByID(id)
	if err != nil {
		fmt.Printf("Error deleting Schema with id=<%d>: %s\n", id, err)
	}

	fmt.Println("END Schema.DeleteByID handler")
	return err
}
