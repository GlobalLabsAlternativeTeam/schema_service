package main

import (
	"log"
	"server/internal/domain"
	"server/internal/handlers/schema"
	"server/internal/providers/storage"
)

func main() {
	// Create instances of your dependencies (handlers, storage, etc.)
	storageService, err := storage.NewStorage("./data/storage.json")
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	schemaHandler := &schema.Schema{StorageProvider: storageService}

	// Generate some tasks

	task1 := domain.Task{
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

	task3 := domain.Task{
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

	task2 := domain.Task{
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

	schemaHandler.Create("Author1", "Schema1", []domain.Task{task1})
	schemaHandler.Create("Author1", "Schema2", []domain.Task{task1, task2})
	schemaHandler.Create("Author2", "Schema3", []domain.Task{task2})
}
