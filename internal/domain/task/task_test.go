package task_test

import (
	"context"
	"crypto/sha1"
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/task"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/mocks/mock_auth"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/mocks/mock_mail"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/mocks/mock_task_storage"
	"go.uber.org/zap"
)

type unitTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &unitTestSuite{})
}

func (s *unitTestSuite) generateToken(task_id string, login string) string {
	logger, _ := zap.NewProduction()

	hash := sha1.New()
	hash.Write([]byte(task_id + login))
	return fmt.Sprintf("%x", hash.Sum([]byte(config.GetConfig(logger.Sugar()).Token.Salt)))
}

func (s *unitTestSuite) generateApprovalLink(task_id string, token string) string {
	logger, _ := zap.NewProduction()
	return "https://127.0.0.1:" + config.GetConfig(logger.Sugar()).Ports.HttpPort + "/tasks/" + task_id + "/approve?token=" + token + "&decision=approve"
}

func (s *unitTestSuite) generateDeclineLink(task_id string, token string) string {
	logger, _ := zap.NewProduction()
	return "https://127.0.0.1:" + config.GetConfig(logger.Sugar()).Ports.HttpPort + "/tasks/" + task_id + "/approve?token=" + token + "&decision=decline"
}

// Tests: CreateTask

func (s *unitTestSuite) TestCreateTaskOk() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("InsertTask",
		&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 0,
			Status:          models.TaskInProgressStatus,
		}).
		Return(nil)

	mm.On("SendApprovalMail",
		models.MailToApproval{
			Destination:  "MyLogin1",
			ApprovalLink: s.generateApprovalLink(taskId, approvalTokens[0]),
			DeclineLink:  s.generateDeclineLink(taskId, approvalTokens[0]),
		}).
		Return()

	taskIn := &models.Task{
		Logins:         logins,
		Title:          title,
		Description:    description,
		InitiatorLogin: initiatorLogin,
	}
	taskOut, err := srv.CreateTask(
		context.Background(),
		taskIn,
	)

	s.Nil(err, "error must be nil")
	s.NotNil(taskOut, "task cannot be nil")
	s.Equal(taskId, taskOut.ID, "wrong task id")
	s.NotEqual(taskIn, taskOut, "returned task object cannot be one from args")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestCreateTaskDbErr() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("InsertTask",
		&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 0,
			Status:          models.TaskInProgressStatus,
		}).
		Return(fmt.Errorf("some error from DB adapter"))

	_, err = srv.CreateTask(
		context.Background(),
		&models.Task{
			Logins:         logins,
			Title:          title,
			Description:    description,
			InitiatorLogin: initiatorLogin,
		},
	)

	s.NotNil(err, "error cannot be nil due to db error")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestCreateTaskNoLogins() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	_, err = srv.CreateTask(
		context.Background(),
		&models.Task{
			Logins:         logins,
			Title:          title,
			Description:    description,
			InitiatorLogin: initiatorLogin,
		},
	)

	s.NotNil(err, "error cannot be nil due to no logins provided")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestCreateTaskEmptyLogin() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"",
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	_, err = srv.CreateTask(
		context.Background(),
		&models.Task{
			Logins:         logins,
			Title:          title,
			Description:    description,
			InitiatorLogin: initiatorLogin,
		},
	)

	s.NotNil(err, "error cannot be nil due to empty login provided")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

// Tests: UpdateTask

func (s *unitTestSuite) TestUpdateTaskOk() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	newLogins := []string{
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	newApprovalTokens := []string{
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"
	userLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 2,
			Status:          models.TaskInProgressStatus,
		}, nil)
	mts.On("UpdateTask",
		&models.Task{
			ID:              taskId,
			Logins:          newLogins,
			ApprovalTokens:  newApprovalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 0,
			Status:          models.TaskInProgressStatus,
		}).
		Return(nil)

	mm.On("SendResultMail",
		models.ResultMail{
			Destinations: logins,
			TaskID:       taskId,
			Result:       "task was updated",
		}).
		Return()
	mm.On("SendApprovalMail",
		models.MailToApproval{
			Destination:  "MyLogin2",
			ApprovalLink: s.generateApprovalLink(taskId, approvalTokens[1]),
			DeclineLink:  s.generateDeclineLink(taskId, approvalTokens[1]),
		}).
		Return()

	taskIn := &models.Task{
		ID:             taskId,
		Logins:         newLogins,
		Title:          title,
		Description:    description,
		InitiatorLogin: initiatorLogin,
	}
	userIn := &models.User{
		Login: userLogin,
	}
	err = srv.UpdateTask(
		context.Background(),
		taskIn,
		userIn,
	)

	s.Nil(err, "error must be nil")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestUpdateTaskNotInitiator() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	newLogins := []string{
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"
	userLogin := "MyLogin"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 2,
			Status:          models.TaskInProgressStatus,
		}, nil)

	taskIn := &models.Task{
		ID:             taskId,
		Logins:         newLogins,
		Title:          title,
		Description:    description,
		InitiatorLogin: initiatorLogin,
	}
	userIn := &models.User{
		Login: userLogin,
	}
	err = srv.UpdateTask(
		context.Background(),
		taskIn,
		userIn,
	)

	s.NotNil(err, "error cannot be nil: user is not task initiator")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

// Tests: ApproveOrDecline

func (s *unitTestSuite) TestApproveOrDeclineNonLastApproveOk() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 2,
			Status:          models.TaskInProgressStatus,
		}, nil)

	mm.On("SendApprovalMail",
		models.MailToApproval{
			Destination:  "MyLogin4",
			ApprovalLink: s.generateApprovalLink(taskId, approvalTokens[3]),
			DeclineLink:  s.generateDeclineLink(taskId, approvalTokens[3]),
		}).
		Return()

	mts.On("UpdateTask",
		&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 3,
			Status:          models.TaskInProgressStatus,
		}).
		Return(nil)

	err = srv.ApproveOrDecline(
		context.Background(),
		taskId,
		approvalTokens[2],
		"approve",
	)

	s.Nil(err, "error must be nil")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestApproveOrDeclineLastApproveOk() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 3,
			Status:          models.TaskInProgressStatus,
		}, nil)

	mm.On("SendResultMail",
		models.ResultMail{
			Destinations: logins,
			TaskID:       taskId,
			Result:       "task was done",
		}).
		Return()

	mts.On("UpdateTask",
		&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 4,
			Status:          models.TaskDoneStatus,
		}).
		Return(nil)

	err = srv.ApproveOrDecline(
		context.Background(),
		taskId,
		approvalTokens[3],
		"approve",
	)

	s.Nil(err, "error must be nil")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestApproveOrDeclineDeclineOk() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 2,
			Status:          models.TaskInProgressStatus,
		}, nil)

	mm.On("SendResultMail",
		models.ResultMail{
			Destinations: logins,
			TaskID:       taskId,
			Result:       "task was cancelled",
		}).
		Return()

	mts.On("UpdateTask",
		&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 2,
			Status:          models.TaskDeclinedStatus,
		}).
		Return(nil)

	err = srv.ApproveOrDecline(
		context.Background(),
		taskId,
		approvalTokens[2],
		"decline",
	)

	s.Nil(err, "error must be nil")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestApproveOrDeclineGetTaskErr() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(nil, fmt.Errorf("failed to get task"))

	err = srv.ApproveOrDecline(
		context.Background(),
		taskId,
		approvalTokens[2],
		"approve",
	)

	s.NotNil(err, "error cannot be nil: get task failed")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestApproveOrDeclineWrongTokenErr() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"
	wrongToken := "1234567890"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 2,
			Status:          models.TaskInProgressStatus,
		}, nil)

	err = srv.ApproveOrDecline(
		context.Background(),
		taskId,
		wrongToken,
		"approve",
	)

	s.NotNil(err, "error cannot be nil: wrong token")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestApproveOrDeclineInvalidDecisionOk() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 2,
			Status:          models.TaskInProgressStatus,
		}, nil)

	err = srv.ApproveOrDecline(
		context.Background(),
		taskId,
		approvalTokens[2],
		"approv",
	)

	s.NotNil(err, "error cannot be nil: invalid decision string")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}

func (s *unitTestSuite) TestApproveOrDeclineUpdateTaskErr() {
	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
	logins := []string{
		"MyLogin1",
		"MyLogin2",
		"MyLogin3",
		"MyLogin4",
	}
	approvalTokens := []string{
		s.generateToken(taskId, "MyLogin1"),
		s.generateToken(taskId, "MyLogin2"),
		s.generateToken(taskId, "MyLogin3"),
		s.generateToken(taskId, "MyLogin4"),
	}
	title := "MyTask1"
	description := "This task is important one too!"
	initiatorLogin := "test123"

	mts := new(mock_task_storage.MockTaskStorage)
	mm := new(mock_mail.MockMail)
	ma := new(mock_auth.MockAuth)
	logger, _ := zap.NewProduction()
	taskIdUuid, err := uuid.FromString(taskId)
	s.NoError(err, "bad task id provided for test")
	srv := task.New(mts, mm, ma, logger.Sugar(), func() uuid.UUID {
		return taskIdUuid
	})

	mts.On("GetTask", taskId).
		Return(&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 2,
			Status:          models.TaskInProgressStatus,
		}, nil)

	mm.On("SendApprovalMail",
		models.MailToApproval{
			Destination:  "MyLogin4",
			ApprovalLink: s.generateApprovalLink(taskId, approvalTokens[3]),
			DeclineLink:  s.generateDeclineLink(taskId, approvalTokens[3]),
		}).
		Return()

	mts.On("UpdateTask",
		&models.Task{
			ID:              taskId,
			Logins:          logins,
			ApprovalTokens:  approvalTokens,
			Title:           title,
			Description:     description,
			InitiatorLogin:  initiatorLogin,
			CurrApprovalNum: 3,
			Status:          models.TaskInProgressStatus,
		}).
		Return(fmt.Errorf("failed to update task"))

	err = srv.ApproveOrDecline(
		context.Background(),
		taskId,
		approvalTokens[2],
		"approve",
	)

	s.NotNil(err, "error cannot be nil: failed to record changes to DB")

	mts.AssertExpectations(s.T())
	mm.AssertExpectations(s.T())
	ma.AssertExpectations(s.T())
}
