package mock_task

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type MockTask struct {
	mock.Mock
}

func (mt *MockTask) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	args := mt.Called(task)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.Task), args.Error(1)
}

func (mt *MockTask) GetTask(ctx context.Context, task_id string) (*models.Task, error) {
	args := mt.Called(task_id)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.Task), args.Error(1)
}

func (mt *MockTask) UpdateTask(ctx context.Context, task *models.Task, user *models.User) error {
	args := mt.Called(task, user)

	return args.Error(0)
}

func (mt *MockTask) DeleteTask(ctx context.Context, task_id string, user *models.User) error {
	args := mt.Called(task_id, user)

	return args.Error(0)
}

func (mt *MockTask) ApproveOrDecline(ctx context.Context, task_id string, token string, decision string) error {
	args := mt.Called(task_id, token, decision)

	return args.Error(0)
}

func (mt *MockTask) ValidateAuth(ctx context.Context, tokenpair *models.TokenPair) (*models.AuthResult, error) {
	args := mt.Called(tokenpair)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.AuthResult), args.Error(1)
}
