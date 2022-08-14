package kafka_producer

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *zap.SugaredLogger
}

func New(brokers []string, topic string, logger *zap.SugaredLogger) (*KafkaProducer, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" {
		logger.Errorf("invalid config parameters for kafka producer")
		return nil, fmt.Errorf("invalid config parameters for kafka producer")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers[0]),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{
		writer: writer,
		logger: logger,
	}, nil
}

func (kp *KafkaProducer) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return kp.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (kp *KafkaProducer) SendMessages(ctx context.Context, messages []*models.Message) error {
	logger := kp.annotatedLogger(ctx)

	for _, m := range messages {
		err := kp.writer.WriteMessages(ctx, kafka.Message{
			Key:   m.Key,
			Value: m.Value,
		})
		if err != nil {
			logger.Errorf("failed to send messages to kafka: %s", err.Error())
			return fmt.Errorf("failed to send messages to kafka: %s", err.Error())
		}
	}

	return nil
}
