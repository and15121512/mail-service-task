package mock_auth

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type MockAuth struct {
	mock.Mock
}

func (ma *MockAuth) ValidateAuth(ctx context.Context, tokenpair *models.TokenPair) (*models.AuthResult, error) {
	args := ma.Called(tokenpair)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.AuthResult), args.Error(1)
}
