package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	filePath string
	schemas  map[string]domain.Schema
}

func NewStorage(filePath string) (*Storage, error) {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If the file doesn't exist, create an empty JSON file
		if err := createEmptyJSONFile(filePath); err != nil {
			return nil, fmt.Errorf("error creating storage file: %v", err)
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading storage file: %v", err)
	}

	var schemas []domain.Schema
	err = json.Unmarshal(data, &schemas)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	schemasMap := make(map[string]domain.Schema)
	for _, schema := range schemas {
		schemasMap[schema.SchemaID] = schema
	}

	return &Storage{
		filePath: filePath,
		schemas:  schemasMap,
	}, nil
}

func (s *Storage) SaveToFile() error {
	schemasSlice := make([]domain.Schema, 0, len(s.schemas))
	for _, schema := range s.schemas {
		schemasSlice = append(schemasSlice, schema)
	}

	data, err := json.MarshalIndent(schemasSlice, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling schemas to JSON: %v", err)
	}

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing data to file: %v", err)
	}

	return nil
}

func createEmptyJSONFile(filePath string) error {
	emptyData := []byte("[]") // JSON representation for empty array
	err := os.WriteFile(filePath, emptyData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateSchema(authorID string, schemaName string, tasks []domain.Task) (domain.Schema, error) {
	fmt.Println("START Storage.CreateSchema")

	// Check if SchemaName is already used
	for _, existingSchema := range s.schemas {
		if existingSchema.SchemaName == schemaName {
			return domain.Schema{}, fmt.Errorf("schema with name '%s' already exists", schemaName)
		}
	}

	// Generate SchemaID
	id := uuid.New().String()
	for { // to avoid (really improbable) collisions
		if _, ok := s.schemas[id]; !ok {
			break
		}
		id = uuid.New().String()
	}

	// Create Schema
	schema := domain.Schema{
		SchemaID:   id,
		AuthorID:   authorID,
		SchemaName: schemaName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Tasks:      tasks,
	}

	// Store in the storage
	s.schemas[id] = schema

	// Save database
	err := s.SaveToFile()
	if err != nil {
		delete(s.schemas, id) // revert changes to avoid broken state
		log.Fatalf("error saving storage to file: %v", err)
		return domain.Schema{}, fmt.Errorf("internal error while creation")
	}

	fmt.Println("END Storage.CreateSchema")
	return schema, nil
}

func (s *Storage) GetSchemaByID(id string) (domain.Schema, error) {
	fmt.Println("START Storage.GetSchemaByID")

	// Get schema and check existance
	schema, ok := s.schemas[id]
	if !ok {
		return domain.Schema{}, fmt.Errorf("schema with id=<%s> not found", id)
	}

	fmt.Println("END Storage.GetSchemaByID")
	return schema, nil

}

func (s *Storage) DeleteSchemaByID(id string) error {
	fmt.Println("START Storage.DeleteSchemaByID")

	// Get schema and check existance
	schema, ok := s.schemas[id]
	if !ok {
		return fmt.Errorf("schema with id=<%s> not found", id)
	}

	// Delete schema from storage
	delete(s.schemas, id)

	// Save database
	err := s.SaveToFile()
	if err != nil {
		s.schemas[id] = schema // revert changes to avoid broken state
		log.Fatalf("error saving storage to file: %v", err)
		return fmt.Errorf("internal error while deletion")
	}

	fmt.Println("END Storage.DeleteSchemaByID")
	return nil

}
