package mail

import (
	"context"
	"encoding/json"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
	"go.uber.org/zap"
)

type Mail struct {
	logger *zap.SugaredLogger
}

func (m *Mail) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return m.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func New(logger *zap.SugaredLogger) *Mail {
	return &Mail{
		logger: logger,
	}
}

func (m *Mail) SendApprovalMail(ctx context.Context, mail models.MailToApproval) {
	logger := m.annotatedLogger(ctx)
	data, _ := json.MarshalIndent(mail, "", "\t")
	logger.Infof("SendApprovalMail: %s", string(data))
}

func (m *Mail) SendResultMail(ctx context.Context, mail models.ResultMail) {
	logger := m.annotatedLogger(ctx)
	data, _ := json.MarshalIndent(mail, "", "\t")
	logger.Infof("SendResultMail: %s", string(data))
}
