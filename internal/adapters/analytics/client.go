package analytics

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
	"go.uber.org/zap"
)

type Analytics struct {
	logger *zap.SugaredLogger
	pr     Producer
}

func New(logger *zap.SugaredLogger, pr Producer) *Analytics {
	return &Analytics{
		logger: logger,
		pr:     pr,
	}
}

func (an *Analytics) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return an.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (an *Analytics) StoreEvent(ctx context.Context, event *models.Event) error {
	logger := an.annotatedLogger(ctx)

	data, err := json.Marshal(event)
	if err != nil {
		logger.Errorf("failed to serialize event structure: %s", err.Error())
		return fmt.Errorf("failed to serialize event structure: %s", err.Error())
	}

	err = an.pr.SendMessages(ctx, []*models.Message{
		&models.Message{
			Key:   []byte(event.EventId),
			Value: data,
		},
	})
	if err != nil {
		logger.Errorf("failed to send event to kafka: %s", err.Error())
		return fmt.Errorf("failed to send event to kafka: %s", err.Error())
	}
	return nil
}
