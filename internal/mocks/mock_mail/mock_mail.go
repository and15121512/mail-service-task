package mock_mail

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
)

type MockMail struct {
	mock.Mock
}

func (mm *MockMail) SendApprovalMail(ctx context.Context, mail models.MailToApproval) {
	_ = mm.Called(mail)
}

func (mm *MockMail) SendResultMail(ctx context.Context, mail models.ResultMail) {
	_ = mm.Called(mail)
}
