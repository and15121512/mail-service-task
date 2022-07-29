package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type Task interface {
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
	GetTask(ctx context.Context, task_id string) (models.Task, error)
	UpdateTask(ctx context.Context, task models.Task, user models.User) error
	DeleteTask(ctx context.Context, task_id string, user models.User) error
	ApproveOrDecline(ctx context.Context, task_id string, token string, decision string) error
	// ...
}
