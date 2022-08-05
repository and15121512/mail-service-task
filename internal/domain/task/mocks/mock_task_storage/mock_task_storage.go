package mock_task_storage

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type MockTaskStorage struct {
	mock.Mock
}

func (mt *MockTaskStorage) InsertTask(ctx context.Context, task *models.Task) error {
	args := mt.Called(task)

	return args.Error(0)
}

func (mt *MockTaskStorage) GetTask(ctx context.Context, task_id string) (*models.Task, error) {
	args := mt.Called(task_id)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.Task), args.Error(1)
}

func (mt *MockTaskStorage) UpdateTask(ctx context.Context, task *models.Task) error {
	args := mt.Called(task)

	return args.Error(0)
}

func (mt *MockTaskStorage) DeleteTask(ctx context.Context, task_id string) error {
	args := mt.Called(task_id)

	return args.Error(0)
}
