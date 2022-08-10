package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type Analytics interface {
	StoreEvent(ctx context.Context, event *models.Event) error
}
