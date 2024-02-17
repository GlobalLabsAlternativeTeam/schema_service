// internal/domain/converter.go

package domain

import (
	"time"

	schema_service "server/proto"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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
		Status:      convertTaskStatus(t.Status),
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
func convertTaskStatus(status string) schema_service.TaskStatus {
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
