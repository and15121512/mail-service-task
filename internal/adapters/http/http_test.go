package http

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type mockTask struct {
	mock.Mock
}

func (mt *mockTask) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	args := mt.Called(task)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.Task), args.Error(1)
}

func (mt *mockTask) GetTask(ctx context.Context, task_id string) (*models.Task, error) {
	args := mt.Called(task_id)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.Task), nil
}

func (mt *mockTask) UpdateTask(ctx context.Context, task *models.Task, user *models.User) error {
	args := mt.Called(task, user)

	return args.Error(0)
}

func (mt *mockTask) DeleteTask(ctx context.Context, task_id string, user *models.User) error {
	args := mt.Called(task_id, user)

	return args.Error(0)
}

func (mt *mockTask) ApproveOrDecline(ctx context.Context, task_id string, token string, decision string) error {
	args := mt.Called(task_id, token, decision)

	return args.Error(0)
}

func (mt *mockTask) ValidateAuth(ctx context.Context, tokenpair *models.TokenPair) (*models.AuthResult, error) {
	args := mt.Called(tokenpair)

	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*models.AuthResult), args.Error(1)
}

// Tests

// type unitTestSuite struct {
// 	suite.Suite
// }

// func TestUnitTestSuite(t *testing.T) {
// 	suite.Run(t, &unitTestSuite{})
// }

// func (s *unitTestSuite) TestCreateTask() {
// 	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
// 	logins := []string{
// 		"MyLogin1",
// 		"MyLogin2",
// 		"MyLogin3",
// 		"MyLogin4",
// 	}
// 	title := "MyTask1"
// 	description := "This task is important one too!"
// 	initiatorLogin := "test123"

// 	taskIn := &models.Task{
// 		Logins:         logins,
// 		Title:          title,
// 		Description:    description,
// 		InitiatorLogin: initiatorLogin,
// 	}
// 	taskOut := &models.Task{
// 		ID: taskId,
// 	}

// 	mt := new(mockTask)
// 	mt.On("CreateTask", taskIn).Return(taskOut, nil)

// 	jsonIn := &struct {
// 		Logins      []string `json:"logins"`
// 		Title       string   `json:"title"`
// 		Description string   `json:"description"`
// 	}{
// 		Logins:      logins,
// 		Title:       title,
// 		Description: description,
// 	}
// 	var b bytes.Buffer
// 	_ = json.NewEncoder(&b).Encode(jsonIn)

// 	r := chi.NewRouter()

// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest("POST", "/tasks", &b)

// 	_ = r
// 	_ = w
// 	_ = req
// }

// func (s *unitTestSuite) TestGetTask() {

// }

// func (s *unitTestSuite) TestUpdateTask() {

// }

// func (s *unitTestSuite) TestDeleteTask() {

// }
