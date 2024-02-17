package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"server/internal/domain"
)

type Storage struct {
	filePath string
	schemas  map[int]domain.Schema
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

	schemasMap := make(map[int]domain.Schema)
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
