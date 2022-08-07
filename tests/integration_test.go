package integration_test

import (
	"context"
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/application"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/config"
	"go.uber.org/zap"
)

type integraTestSuite struct {
	suite.Suite
}

func TestIntegraTestSuite(t *testing.T) {
	suite.Run(t, &integraTestSuite{})
}

func (s *integraTestSuite) generateToken(task_id string, login string) string {
	logger, _ := zap.NewProduction()

	hash := sha1.New()
	hash.Write([]byte(task_id + login))
	return fmt.Sprintf("%x", hash.Sum([]byte(config.GetConfig(logger.Sugar()).Token.Salt)))
}

func (s *integraTestSuite) generateApprovalLink(task_id string, token string) string {
	logger, _ := zap.NewProduction()
	return "https://127.0.0.1:" + config.GetConfig(logger.Sugar()).Ports.HttpPort + "/tasks/" + task_id + "/approve?token=" + token + "&decision=approve"
}

func (s *integraTestSuite) generateDeclineLink(task_id string, token string) string {
	logger, _ := zap.NewProduction()
	return "https://127.0.0.1:" + config.GetConfig(logger.Sugar()).Ports.HttpPort + "/tasks/" + task_id + "/approve?token=" + token + "&decision=decline"
}

func (s *integraTestSuite) SetupSuite() {
	go application.Start(context.Background())
}

func (s *integraTestSuite) TearDownSuite() {
	application.Stop()
}
