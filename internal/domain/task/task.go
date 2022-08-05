package task

import (
	"context"
	"crypto/sha1"
	"fmt"
	"regexp"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/ports"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
	"go.uber.org/zap"
)

type Service struct {
	ts       ports.TaskStorage
	m        ports.Mail
	ac       ports.Auth
	logger   *zap.SugaredLogger
	uuidFunc func() uuid.UUID
}

func New(ts ports.TaskStorage, m ports.Mail, ac ports.Auth, logger *zap.SugaredLogger, uuidFunc ...func() uuid.UUID) *Service {
	s := &Service{
		ts:       ts,
		m:        m,
		ac:       ac,
		logger:   logger,
		uuidFunc: uuid.NewV4,
	}

	if len(uuidFunc) > 0 {
		s.uuidFunc = uuidFunc[0]
	}
	return s
}

func (s *Service) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return s.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (s *Service) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	logger := s.annotatedLogger(ctx)
	if err := s.isTaskValid(task); err != nil {
		logger.Errorf(err.Error())
		return &models.Task{}, err
	}

	newTask := *task
	newTask.ID = s.uuidFunc().String()
	approvalTokens := make([]string, len(newTask.Logins))
	for i, login := range newTask.Logins {
		approvalTokens[i] = s.generateToken(ctx, newTask.ID, login)
	}
	newTask.ApprovalTokens = approvalTokens
	newTask.CurrApprovalNum = 0
	newTask.Status = models.TaskInProgressStatus

	err := s.ts.InsertTask(ctx, &newTask)
	if err != nil {
		logger.Errorf("failed to insert task into storage")
		return &models.Task{}, fmt.Errorf("failed to insert task into storage")
	}

	s.m.SendApprovalMail(ctx, models.MailToApproval{
		Destination:  newTask.Logins[0],
		ApprovalLink: s.generateApprovalLink(newTask.ID, approvalTokens[0]),
		DeclineLink:  s.generateDeclineLink(newTask.ID, approvalTokens[0]),
	})

	return &newTask, nil
}

func (s *Service) isTaskValid(task *models.Task) error {
	if len(task.Logins) == 0 {
		return fmt.Errorf("no logins provided to send approval emails")
	}
	for _, login := range task.Logins {
		if err := s.isLoginValid(login); err != nil {
			return fmt.Errorf("invalid email provided: %s", err.Error())
		}
	}
	return nil
}

func (s *Service) isLoginValid(login string) error {
	pattern := `\w+`
	matched, err := regexp.Match(pattern, []byte(login))
	if err != nil || !matched {
		return fmt.Errorf("failed to parse login: only [A-Za-z0-9_]+ allowed")
	}
	return nil
}

func (s *Service) GetTask(ctx context.Context, task_id string) (*models.Task, error) {
	return s.ts.GetTask(ctx, task_id)
}

func (s *Service) UpdateTask(ctx context.Context, task *models.Task, user *models.User) error {
	logger := s.annotatedLogger(ctx)
	if err := s.isTaskValid(task); err != nil {
		logger.Errorf(err.Error())
		return err
	}

	newTask := *task
	oldTask, err := s.ts.GetTask(ctx, newTask.ID)
	if err != nil {
		logger.Errorf("failed to get task with task ID %s for update", newTask.ID)
		return fmt.Errorf("failed to get task with task ID %s for update", newTask.ID)
	}
	if user.Login != oldTask.InitiatorLogin {
		logger.Errorf("user %s is not the task %s author", user.Login, oldTask.ID)
		return fmt.Errorf("user %s is not the task %s author", user.Login, oldTask.ID)
	}

	approvalTokens := make([]string, len(newTask.Logins))
	for i, login := range newTask.Logins {
		approvalTokens[i] = s.generateToken(ctx, newTask.ID, login)
	}
	newTask.ApprovalTokens = approvalTokens
	newTask.InitiatorLogin = oldTask.InitiatorLogin
	newTask.CurrApprovalNum = 0
	newTask.Status = models.TaskInProgressStatus

	err = s.ts.UpdateTask(ctx, &newTask)
	if err != nil {
		logger.Errorf("failed to update task with task ID %s", newTask.ID)
		return fmt.Errorf("failed to update task with task ID %s", newTask.ID)
	}

	s.m.SendResultMail(ctx, models.ResultMail{
		Destinations: oldTask.Logins,
		TaskID:       oldTask.ID,
		Result:       "task was updated",
	})
	s.m.SendApprovalMail(ctx, models.MailToApproval{
		Destination:  newTask.Logins[0],
		ApprovalLink: s.generateApprovalLink(newTask.ID, approvalTokens[0]),
		DeclineLink:  s.generateDeclineLink(newTask.ID, approvalTokens[0]),
	})
	return nil
}

func (s *Service) DeleteTask(ctx context.Context, task_id string, user *models.User) error {
	logger := s.annotatedLogger(ctx)

	task, err := s.ts.GetTask(ctx, task_id)
	if err != nil {
		logger.Errorf("failed to get task with task ID %s for update", task_id)
		return fmt.Errorf("failed to get task with task ID %s for update", task_id)
	}
	if user.Login != task.InitiatorLogin {
		logger.Errorf("user %s is not the task %s author", user.Login, task_id)
		return fmt.Errorf("user %s is not the task %s author", user.Login, task_id)
	}

	err = s.ts.DeleteTask(ctx, task_id)
	if err != nil {
		logger.Errorf("failed to delete task with task ID %s", task_id)
		return fmt.Errorf("failed to delete task with task ID %s", task_id)
	}

	s.m.SendResultMail(ctx, models.ResultMail{
		Destinations: task.Logins,
		TaskID:       task.ID,
		Result:       "task was deleted",
	})

	return nil
}

func (s *Service) ApproveOrDecline(ctx context.Context, task_id string, token string, decision string) error {
	logger := s.annotatedLogger(ctx)

	task, err := s.ts.GetTask(ctx, task_id)
	if err != nil {
		logger.Errorf("failed to get task with task ID %s", task_id)
		return fmt.Errorf("failed to get task with task ID %s", task_id)
	}

	if task.ApprovalTokens[task.CurrApprovalNum] != token {
		logger.Errorf("invalid approval token")
		return fmt.Errorf("invalid approval token")
	}

	if decision == "approve" {
		err = s.approve(ctx, task)
	} else if decision == "decline" {
		err = s.decline(ctx, task)
	} else {
		logger.Errorf("invalid decision value (expected 'approve' or 'decline')")
		return fmt.Errorf("invalid decision value (expected 'approve' or 'decline')")
	}
	if err != nil {
		logger.Errorf("failed to process approve/decline: %s", err.Error())
		return fmt.Errorf("failed to process approve/decline: %s", err.Error())
	}

	err = s.ts.UpdateTask(ctx, task)
	if err != nil {
		logger.Errorf("failed to update task with task ID %s", task.ID)
		return fmt.Errorf("failed to update task with task ID %s", task.ID)
	}

	return nil
}

func (s *Service) approve(ctx context.Context, task *models.Task) error {
	task.CurrApprovalNum++
	if task.CurrApprovalNum >= len(task.Logins) {
		task.Status = models.TaskDoneStatus
		s.m.SendResultMail(ctx, models.ResultMail{
			Destinations: task.Logins,
			TaskID:       task.ID,
			Result:       "task was done",
		})
	} else {
		if len(task.Logins) == 0 {
			return fmt.Errorf("logins array of the task %s is empty", task.ID)
		}
		s.m.SendApprovalMail(ctx, models.MailToApproval{
			Destination:  task.Logins[task.CurrApprovalNum],
			ApprovalLink: s.generateApprovalLink(task.ID, task.ApprovalTokens[task.CurrApprovalNum]),
			DeclineLink:  s.generateDeclineLink(task.ID, task.ApprovalTokens[task.CurrApprovalNum]),
		})
	}
	return nil
}

func (s *Service) decline(ctx context.Context, task *models.Task) error {
	task.Status = models.TaskDeclinedStatus
	s.m.SendResultMail(ctx, models.ResultMail{
		Destinations: task.Logins,
		TaskID:       task.ID,
		Result:       "task was cancelled",
	})
	return nil
}

func (s *Service) generateToken(ctx context.Context, task_id string, login string) string {
	logger := s.annotatedLogger(ctx)

	hash := sha1.New()
	hash.Write([]byte(task_id + login))
	return fmt.Sprintf("%x", hash.Sum([]byte(config.GetConfig(logger).Token.Salt)))
}

func (s *Service) generateApprovalLink(task_id string, token string) string {
	return "https://127.0.0.1:" + config.GetConfig(s.logger).Ports.HttpPort + "/tasks/" + task_id + "/approve?token=" + token + "&decision=approve"
}

func (s *Service) generateDeclineLink(task_id string, token string) string {
	return "https://127.0.0.1:" + config.GetConfig(s.logger).Ports.HttpPort + "/tasks/" + task_id + "/approve?token=" + token + "&decision=decline"
}
