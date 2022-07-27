package application

import (
	"context"
	"fmt"

	"github.com/TheZeroSlave/zapsentry"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/adapters/data_file"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/adapters/http"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/adapters/mail"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/task"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

var (
	s      *http.Server
	logger *zap.Logger
)

func modifyToSentryLogger(log *zap.Logger, DSN string) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(DSN))

	// to use breadcrumbs feature - create new scope explicitly
	log = log.With(zapsentry.NewScope())

	//in case of err it will return noop core. so we can safely attach it
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}
	return zapsentry.AttachCoreToLogger(core, log)
}

func Start(ctx context.Context) {
	logger, _ = zap.NewProduction()
	//sentryClient, err := sentry.NewClient(sentry.ClientOptions{
	//	Dsn: "http://b7dd7b3ce3df4f2b81f5af622512658c@localhost:9000/2",
	//})
	//if err != nil {
	//	logger.Sugar().Fatalf("http server creating failed: %s", err)
	//}
	//defer sentryClient.Flush(2 * time.Second)
	logger = modifyToSentryLogger(logger, "http://b7dd7b3ce3df4f2b81f5af622512658c@localhost:9000/2")

	db, err := data_file.New(logger.Sugar()) // Mock!!!
	if err != nil {
		logger.Sugar().Fatalf("db init failed: %s", err)
	}
	m := mail.New(logger.Sugar())

	taskS := task.New(db, m, logger.Sugar())

	s, err = http.New(taskS, logger.Sugar())
	if err != nil {
		logger.Sugar().Fatalf("http server creating failed: %s", err)
	}

	var g errgroup.Group
	g.Go(func() error {
		return s.Start()
	})

	logger.Sugar().Info(fmt.Sprintf("app is started on port:%d", s.Port()))
	err = g.Wait()
	if err != nil {
		logger.Sugar().Fatalw("http server start failed", zap.Error(err))
	}
}

func Stop() {
	_ = s.Stop(context.Background())
	logger.Sugar().Info("app has stopped")
}