package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type Auth interface {
	ValidateAuth(ctx context.Context, tokenpair models.TokenPair) (models.AuthResult, error)
}
