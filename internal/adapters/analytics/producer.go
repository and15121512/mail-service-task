package analytics

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type Producer interface {
	SendMessages(ctx context.Context, messages []*models.Message) error
}
