package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type Mail interface {
	SendApprovalMail(ctx context.Context, mail models.MailToApproval)
	SendResultMail(ctx context.Context, mail models.ResultMail)
}
