package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type TaskStorage interface {
	InsertTask(ctx context.Context, task models.Task) error
}
