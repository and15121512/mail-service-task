package task

import (
	"context"
	"crypto/sha1"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/ports"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
	"go.uber.org/zap"
)

type Service struct {
	ts     ports.TaskStorage
	mail   ports.Mail
	logger *zap.SugaredLogger
}

func New(ts ports.TaskStorage, mail ports.Mail, logger *zap.SugaredLogger) *Service {
	return &Service{
		ts:     ts,
		mail:   mail,
		logger: logger,
	}
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

func (s *Service) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	logger := s.annotatedLogger(ctx)

	task.ID = uuid.NewV4().String()
	approvalTokens := make([]string, len(task.Logins))
	for i, login := range task.Logins {
		approvalTokens[i] = s.generateToken(ctx, task.ID, login)
	}
	task.ApprovalTokens = approvalTokens
	task.CurrApprovalNum = 0
	task.Status = models.TaskInProgressStatus

	err := s.ts.InsertTask(ctx, task)
	if err != nil {
		logger.Errorf("failed to insert task into storage")
		return models.Task{}, fmt.Errorf("failed to insert task into storage")
	}

	if len(task.Logins) == 0 {
		return task, nil
	}
	s.mail.SendApprovalMail(ctx, models.MailToApproval{
		Destination:  task.Logins[0],
		ApprovalLink: s.generateApprovalLink(task.ID, approvalTokens[0]),
		DeclineLink:  s.generateDeclineLink(task.ID, approvalTokens[0]),
	})

	return task, nil
}

func (s *Service) generateToken(ctx context.Context, task_id string, login string) string {
	logger := s.annotatedLogger(ctx)

	hash := sha1.New()
	hash.Write([]byte(task_id + login))
	return fmt.Sprintf("%x", hash.Sum([]byte(config.GetConfig(logger).Token.Salt)))
}

func (s *Service) generateApprovalLink(task_id string, token string) string {
	return "https://127.0.0.1:" + config.GetConfig(s.logger).Listen.Port + "/tasks/" + task_id + "/approve?token=" + token + "&decision=approve"
}

func (s *Service) generateDeclineLink(task_id string, token string) string {
	return "https://127.0.0.1:" + config.GetConfig(s.logger).Listen.Port + "/tasks/" + task_id + "/approve?token=" + token + "&decision=decline"
}
