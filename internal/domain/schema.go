package domain

import "time"

type Task struct {
	ID          int     `json:"id"`
	Level       int     `json:"level"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	BlockedBy   []int64 `json:"blocked_by"`
	Responsible string  `json:"responsible"`
	TimeLimit   int64   `json:"time_limit"`
	Children    []Task  `json:"children"`
	Comment     struct {
		Value string `json:"value"`
	} `json:"comment"`
}

type Schema struct {
	SchemaID   int       `json:"schema_id"`
	AuthorID   string    `json:"author_id"`
	SchemaName string    `json:"schema_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
	Tasks      []Task    `json:"tasks"`
}
