package analytics_test

import (
	"context"
	"testing"
	"time"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/adapters/analytics"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/adapters/analytics/kafka_producer"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"go.uber.org/zap"
)

func TestStoreEvent(t *testing.T) {
	logger, _ := zap.NewProduction()

	kp, err := kafka_producer.New(
		[]string{"localhost:9092"},
		"task_event",
		logger.Sugar(),
	)
	if err != nil {
		t.Fatalf("failed to create kafka producer: %s", err.Error())
	}

	an := analytics.New(logger.Sugar(), kp)
	err = an.StoreEvent(
		context.Background(),
		&models.Event{
			EventId: "d8e84c49-4d63-4cba-a653-955b5e35001b",
			TaskId:  "d8e84c49-4d63-4cba-a653-955b5e35001b",
			Time:    time.Now(),
			Type:    models.EventCreateType,
			Status:  models.TaskInProgressStatus,
		},
	)
	if err != nil {
		t.Fatalf("failed to store event in kafka: %s", err.Error())
	}
}
