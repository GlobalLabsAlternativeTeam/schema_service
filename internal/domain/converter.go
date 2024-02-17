// internal/domain/converter.go

package domain

import (
	"time"

	schema_service "server/proto"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// To gRPC

func SchemaToGRPC(s *Schema) *schema_service.Schema {
	return &schema_service.Schema{
		SchemaId:   s.SchemaID,
		AuthorId:   s.AuthorID,
		SchemaName: s.SchemaName,
		CreatedAt:  convertTimestampFromTime(s.CreatedAt),
		UpdatedAt:  convertTimestampFromTime(s.UpdatedAt),
		DeletedAt:  convertTimestampFromTime(s.DeletedAt),
		Tasks:      TasksToGRPC(s.Tasks),
	}
}

func TasksToGRPC(tasks []Task) []*schema_service.Task {
	var grpcTasks []*schema_service.Task
	for _, t := range tasks {
		grpcTask := TaskToGRPC(&t)
		grpcTasks = append(grpcTasks, grpcTask)
	}
	return grpcTasks
}

func TaskToGRPC(t *Task) *schema_service.Task {
	return &schema_service.Task{
		Id:          int64(t.ID),
		Level:       int32(t.Level),
		Name:        t.Name,
		Status:      convertTaskStatusToGRPC(t.Status),
		BlockedBy:   t.BlockedBy,
		Responsible: t.Responsible,
		TimeLimit:   t.TimeLimit,
		Children:    TasksToGRPC(t.Children),
		Comment:     wrapperspb.String(t.Comment.Value),
	}
}

func convertTimestampFromTime(t time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: int64(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}
}

func convertTaskStatusToGRPC(status string) schema_service.TaskStatus {
	switch status {
	case "NOT_STARTED":
		return schema_service.TaskStatus_TASK_STATUS_NOT_STARTED
	case "IN_PROGRESS":
		return schema_service.TaskStatus_TASK_STATUS_IN_PROGRESS
	case "BLOCKED":
		return schema_service.TaskStatus_TASK_STATUS_BLOCKED
	case "DONE":
		return schema_service.TaskStatus_TASK_STATUS_DONE
	default:
		return schema_service.TaskStatus_TASK_STATUS_UNSPECIFIED
	}
}

// From gRPC

func TasksFromGRPC(grpcTasks []*schema_service.Task) []Task {
	var tasks []Task
	for _, grpcTask := range grpcTasks {
		task := TaskFromGRPC(grpcTask)
		tasks = append(tasks, task)
	}

	return tasks
}

func TaskFromGRPC(t *schema_service.Task) Task {
	return Task{
		ID:          int(t.Id),
		Level:       int(t.Level),
		Name:        t.Name,
		Status:      convertTaskStatusFromGRPC(&t.Status),
		BlockedBy:   t.BlockedBy,
		Responsible: t.Responsible,
		TimeLimit:   t.TimeLimit,
		Children:    TasksFromGRPC(t.Children),
		Comment: struct {
			Value string "json:\"value\""
		}{
			Value: t.Comment.Value,
		},
	}
}

func convertTaskStatusFromGRPC(status *schema_service.TaskStatus) string {
	switch *status {
	case schema_service.TaskStatus_TASK_STATUS_NOT_STARTED:
		return "NOT_STARTED"
	case schema_service.TaskStatus_TASK_STATUS_IN_PROGRESS:
		return "IN_PROGRESS"
	case schema_service.TaskStatus_TASK_STATUS_BLOCKED:
		return "BLOCKED"
	case schema_service.TaskStatus_TASK_STATUS_DONE:
		return "DONE"
	default:
		return "UNSPECIFIED"
	}
}
