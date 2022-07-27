package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type Task interface {
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
	// ...
}
