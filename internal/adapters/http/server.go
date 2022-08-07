package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/ports"
	"go.uber.org/zap"
)

type Server struct {
	task          ports.Task
	server        *http.Server
	l             net.Listener
	port          int
	logger        *zap.SugaredLogger
	defaultTaskId string
}

func New(task ports.Task, logger *zap.SugaredLogger) (*Server, error) {
	var (
		err error
		s   Server
	)
	s.logger = logger
	s.l, err = net.Listen("tcp", ":"+config.GetConfig(logger).Ports.HttpPort)
	if err != nil {
		logger.Fatalf("Failed listen port: %s", err)
	}
	s.task = task
	s.port = s.l.Addr().(*net.TCPAddr).Port
	s.server = &http.Server{
		Handler: s.routes(),
	}

	return &s, nil
}

func NewTest(task ports.Task, logger *zap.SugaredLogger, defaultTaskId ...string) *Server {
	var (
		s Server
	)
	s.logger = logger
	s.task = task

	s.defaultTaskId = ""
	if len(defaultTaskId) > 0 {
		s.defaultTaskId = defaultTaskId[0]
	}

	return &s
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Start() error {
	if err := s.server.Serve(s.l); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/", s.taskHandlers())
	r.Mount("/debug/", middleware.Profiler())

	return r
}
