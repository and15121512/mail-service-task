package analytics_grpc

import (
	"context"
	"fmt"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
	"gitlab.com/sukharnikov.aa/mail-service-task/pkg/analyticsgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalyticsGrpc struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) *AnalyticsGrpc {
	return &AnalyticsGrpc{
		logger: logger,
	}
}

func (ag *AnalyticsGrpc) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return ag.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (ag *AnalyticsGrpc) StoreEvent(ctx context.Context, event *models.Event) error {
	logger := ag.annotatedLogger(ctx)

	analytics_host := config.GetConfig(logger).Hosts.AnalyticsHost
	analytics_port := config.GetConfig(logger).Ports.GrpcPort
	target := fmt.Sprintf("%s:%s", analytics_host, analytics_port)
	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		logger.Errorf("failed to connect to analytics service: %s", err.Error())
		return fmt.Errorf("failed to connect to analytics service: %s", err.Error())
	}
	defer conn.Close()

	a := analyticsgrpc.NewAnalyticsClient(conn)

	analyticsEvent := analyticsgrpc.Event{
		EventId: event.EventId,
		TaskId:  event.TaskId,
		Time:    timestamppb.New(event.Time),
	}
	switch event.Type {
	case models.EventCreateType:
		analyticsEvent.Type = "create"
	case models.EventUpdateType:
		analyticsEvent.Type = "update"
	case models.EventDeleteType:
		analyticsEvent.Type = "delete"
	case models.EventApproveType:
		analyticsEvent.Type = "approve"
	case models.EventDeclineType:
		analyticsEvent.Type = "decline"
	default:
		return fmt.Errorf("unknown event type code")
	}
	switch event.Status {
	case models.TaskInProgressStatus:
		analyticsEvent.Status = "in_progress"
	case models.TaskDoneStatus:
		analyticsEvent.Status = "done"
	case models.TaskDeclinedStatus:
		analyticsEvent.Status = "declined"
	default:
		return fmt.Errorf("unknown task status code")
	}

	_, err = a.StoreEvent(ctx, &analyticsEvent)
	if err != nil {
		logger.Errorf("failed to call store event from analytics service: %s", err.Error())
		return fmt.Errorf("failed to call store event from analytics service: %s", err.Error())
	}
	return nil
}
