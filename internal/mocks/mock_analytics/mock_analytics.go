package mock_analytics

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type MockAnalytics struct {
	mock.Mock
}

func (ma *MockAnalytics) StoreEvent(ctx context.Context, event *models.Event) error {
	args := ma.Called() // 'event' should be added

	return args.Error(0)
}
