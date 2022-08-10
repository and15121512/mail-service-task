package http_test

// import (
// 	"bytes"
// 	"crypto/sha1"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/suite"
// 	task_http "gitlab.com/sukharnikov.aa/mail-service-task/internal/adapters/http"
// 	"gitlab.com/sukharnikov.aa/mail-service-task/internal/adapters/mail"
// 	"gitlab.com/sukharnikov.aa/mail-service-task/internal/config"
// 	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
// 	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/task"
// 	"gitlab.com/sukharnikov.aa/mail-service-task/internal/mocks/mock_auth"
// 	"gitlab.com/sukharnikov.aa/mail-service-task/internal/mocks/mock_task"
// 	"gitlab.com/sukharnikov.aa/mail-service-task/internal/mocks/mock_task_storage"
// 	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
// 	"go.uber.org/zap"
// )

// const (
// 	testAccessToken           = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTk4MTY3ODEsImlhdCI6MTY1OTgxNjcyMSwibG9naW4iOiJ0ZXN0MTIzIn0.0yML4R54Bp7AYOjFZ61mJCCQuRREsDiJxrxg_QLXK2E"
// 	testRefreshToken          = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTk4MjAzMjEsImlhdCI6MTY1OTgxNjcyMSwibG9naW4iOiJ0ZXN0MTIzIn0.-VUgsoyxScxXSwKbvc5qFshog50zqrwqtJoZ2-aIS2E"
// 	testAccessTokenRefreshed  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTk4NjYyNjMsImlhdCI6MTY1OTg2NjIwMywibG9naW4iOiJ0ZXN0MTIzIn0.xFaBWShd_f2HoZs7nunEoI8S_y9Wt3kGSf3tcnlsfEg"
// 	testRefreshTokenRefreshed = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTk4Njk4MDMsImlhdCI6MTY1OTg2NjIwMywibG9naW4iOiJ0ZXN0MTIzIn0.xxqN3lPY766IJT54OsSloCzy4r8WQ1bdgiHAm_CZowg"
// )

// type unitTestSuite struct {
// 	suite.Suite
// }

// func TestUnitTestSuite(t *testing.T) {
// 	suite.Run(t, &unitTestSuite{})
// }

// func (s *unitTestSuite) generateToken(task_id string, login string) string {
// 	logger, _ := zap.NewProduction()

// 	hash := sha1.New()
// 	hash.Write([]byte(task_id + login))
// 	return fmt.Sprintf("%x", hash.Sum([]byte(config.GetConfig(logger.Sugar()).Token.Salt)))
// }

// // Tests: CreateTask

// func (s *unitTestSuite) TestValidateAuthNonEmptyRefreshTokenOk() {
// 	mt := new(mock_task.MockTask)
// 	logger, _ := zap.NewProduction()
// 	srv := task_http.NewTest(mt, logger.Sugar())

// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		userLogin := r.Context().Value(utils.CtxKeyUserGet())
// 		s.NotNil(userLogin, "no user login key in context")
// 		reqId := r.Context().Value(utils.CtxKeyRequestIDGet())
// 		s.NotNil(reqId, "no request ID key in context")
// 		method := r.Context().Value(utils.CtxKeyMethodGet())
// 		s.NotNil(method, "no method key in context")
// 		url := r.Context().Value(utils.CtxKeyURLGet())
// 		s.NotNil(url, "no URL key in context")
// 	})
// 	validateMiddleware := srv.AnnotateContext()(srv.ValidateAuth()(nextHandler))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	mt.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        "test123",
// 		}, nil)

// 	validateMiddleware.ServeHTTP(rec, req)
// }

// func (s *unitTestSuite) TestValidateAuthEmptyRefreshTokenOk() {
// 	mt := new(mock_task.MockTask)
// 	logger, _ := zap.NewProduction()
// 	srv := task_http.NewTest(mt, logger.Sugar())

// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		userLogin := r.Context().Value(utils.CtxKeyUserGet())
// 		s.NotNil(userLogin, "no user login key in context")
// 		reqId := r.Context().Value(utils.CtxKeyRequestIDGet())
// 		s.NotNil(reqId, "no request ID key in context")
// 		method := r.Context().Value(utils.CtxKeyMethodGet())
// 		s.NotNil(method, "no method key in context")
// 		url := r.Context().Value(utils.CtxKeyURLGet())
// 		s.NotNil(url, "no URL key in context")
// 	})
// 	validateMiddleware := srv.AnnotateContext()(srv.ValidateAuth()(nextHandler))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})

// 	mt.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: "",
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        "test123",
// 		}, nil)

// 	validateMiddleware.ServeHTTP(rec, req)
// }

// func (s *unitTestSuite) TestValidateAuthRefreshedOk() {
// 	mt := new(mock_task.MockTask)
// 	logger, _ := zap.NewProduction()
// 	srv := task_http.NewTest(mt, logger.Sugar())

// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		userLogin := r.Context().Value(utils.CtxKeyUserGet())
// 		s.NotNil(userLogin, "no user login key in context")
// 		reqId := r.Context().Value(utils.CtxKeyRequestIDGet())
// 		s.NotNil(reqId, "no request ID key in context")
// 		method := r.Context().Value(utils.CtxKeyMethodGet())
// 		s.NotNil(method, "no method key in context")
// 		url := r.Context().Value(utils.CtxKeyURLGet())
// 		s.NotNil(url, "no URL key in context")
// 	})
// 	validateMiddleware := srv.AnnotateContext()(srv.ValidateAuth()(nextHandler))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	mt.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.RefreshedAuthRespStatus,
// 			AccessToken:  testAccessTokenRefreshed,
// 			RefreshToken: testRefreshTokenRefreshed,
// 			Login:        "test123",
// 		}, nil)

// 	validateMiddleware.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	accessVal, err := getCookieFromResp(resp, "access")
// 	s.Nil(err)
// 	s.Equal(testAccessTokenRefreshed, accessVal)
// 	refreshVal, err := getCookieFromResp(resp, "refresh")
// 	s.Nil(err)
// 	s.Equal(testRefreshTokenRefreshed, refreshVal)
// }

// func getCookieFromResp(resp *http.Response, name string) (string, error) {
// 	for _, cookie := range resp.Cookies() {
// 		if cookie.Name == name {
// 			return cookie.Value, nil
// 		}
// 	}
// 	return "", fmt.Errorf("cookie %s not found in resp", name)
// }

// func (s *unitTestSuite) TestValidateAuthEmptyAccessTokenErr() {
// 	mt := new(mock_task.MockTask)
// 	logger, _ := zap.NewProduction()
// 	srv := task_http.NewTest(mt, logger.Sugar())

// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	})
// 	validateMiddleware := srv.AnnotateContext()(srv.ValidateAuth()(nextHandler))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", nil)

// 	validateMiddleware.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"error":"cannot read token from header"}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusForbidden, resp.StatusCode)
// }

// func (s *unitTestSuite) TestValidateAuthAuthReturnedErrorErr() {
// 	mt := new(mock_task.MockTask)
// 	logger, _ := zap.NewProduction()
// 	srv := task_http.NewTest(mt, logger.Sugar())

// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	})
// 	validateMiddleware := srv.AnnotateContext()(srv.ValidateAuth()(nextHandler))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	mt.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(nil, fmt.Errorf("auth failed to process given tokens"))

// 	validateMiddleware.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"error":"failed to validate token with auth service: auth failed to process given tokens"}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusInternalServerError, resp.StatusCode)
// }

// func (s *unitTestSuite) TestValidateAuthRefusedErr() {
// 	mt := new(mock_task.MockTask)
// 	logger, _ := zap.NewProduction()
// 	srv := task_http.NewTest(mt, logger.Sugar())

// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	})
// 	validateMiddleware := srv.AnnotateContext()(srv.ValidateAuth()(nextHandler))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	mt.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.RefusedAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        "",
// 		}, nil)

// 	validateMiddleware.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"error":"authorization required"}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusUnauthorized, resp.StatusCode)
// }

// // Create Task

// func (s *unitTestSuite) TestCreateTaskOk() {
// 	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
// 	logins := []string{
// 		"MyLogin1",
// 		"MyLogin2",
// 		"MyLogin3",
// 		"MyLogin4",
// 	}
// 	approvalTokens := []string{
// 		s.generateToken(taskId, "MyLogin1"),
// 		s.generateToken(taskId, "MyLogin2"),
// 		s.generateToken(taskId, "MyLogin3"),
// 		s.generateToken(taskId, "MyLogin4"),
// 	}
// 	title := "MyTask1"
// 	description := "This task is important one too!"
// 	initiatorLogin := "test123"

// 	logger, _ := zap.NewProduction()
// 	mts := new(mock_task_storage.MockTaskStorage)
// 	ma := new(mock_auth.MockAuth)
// 	m := mail.New(logger.Sugar())
// 	t := task.New(mts, m, ma, logger.Sugar(), func() uuid.UUID {
// 		taskIdUuid, _ := uuid.FromString(taskId)
// 		return taskIdUuid
// 	})
// 	srv := task_http.NewTest(t, logger.Sugar())

// 	createTaskHandler := srv.AnnotateContext()(srv.ValidateAuth()(http.HandlerFunc(srv.CreateTask)))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(`{
// 		"logins": [ "MyLogin1", "MyLogin2", "MyLogin3", "MyLogin4" ],
// 		"title": "MyTask1",
// 		"description": "This task is important one too!"
// 	}`)))

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	ma.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        initiatorLogin,
// 		}, nil)

// 	mts.On("InsertTask",
// 		&models.Task{
// 			ID:              taskId,
// 			Logins:          logins,
// 			ApprovalTokens:  approvalTokens,
// 			Title:           title,
// 			Description:     description,
// 			InitiatorLogin:  initiatorLogin,
// 			CurrApprovalNum: 0,
// 			Status:          models.TaskInProgressStatus,
// 		}).
// 		Return(nil)

// 	createTaskHandler.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"task_id":"1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusCreated, resp.StatusCode)
// }

// func (s *unitTestSuite) TestCreateTaskBadRequestErr() {
// 	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
// 	initiatorLogin := "test123"

// 	logger, _ := zap.NewProduction()
// 	mts := new(mock_task_storage.MockTaskStorage)
// 	ma := new(mock_auth.MockAuth)
// 	m := mail.New(logger.Sugar())
// 	t := task.New(mts, m, ma, logger.Sugar(), func() uuid.UUID {
// 		taskIdUuid, _ := uuid.FromString(taskId)
// 		return taskIdUuid
// 	})
// 	srv := task_http.NewTest(t, logger.Sugar())

// 	createTaskHandler := srv.AnnotateContext()(srv.ValidateAuth()(http.HandlerFunc(srv.CreateTask)))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(`{
// 		"logins": [ "MyLogin1", "MyLogin2", "MyLogin3", "MyLogin4" ],
// 		"title": "MyTask1",
// 		"description": "This task is important one too!"
// 	`)))

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	ma.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        initiatorLogin,
// 		}, nil)

// 	createTaskHandler.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"error":"cannot parse input for create task request"}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusBadRequest, resp.StatusCode)
// }

// func (s *unitTestSuite) TestCreateTaskCreateFailedErr() {
// 	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
// 	logins := []string{
// 		"MyLogin1",
// 		"MyLogin2",
// 		"MyLogin3",
// 		"MyLogin4",
// 	}
// 	approvalTokens := []string{
// 		s.generateToken(taskId, "MyLogin1"),
// 		s.generateToken(taskId, "MyLogin2"),
// 		s.generateToken(taskId, "MyLogin3"),
// 		s.generateToken(taskId, "MyLogin4"),
// 	}
// 	title := "MyTask1"
// 	description := "This task is important one too!"
// 	initiatorLogin := "test123"

// 	logger, _ := zap.NewProduction()
// 	mts := new(mock_task_storage.MockTaskStorage)
// 	ma := new(mock_auth.MockAuth)
// 	m := mail.New(logger.Sugar())
// 	t := task.New(mts, m, ma, logger.Sugar(), func() uuid.UUID {
// 		taskIdUuid, _ := uuid.FromString(taskId)
// 		return taskIdUuid
// 	})
// 	srv := task_http.NewTest(t, logger.Sugar())

// 	createTaskHandler := srv.AnnotateContext()(srv.ValidateAuth()(http.HandlerFunc(srv.CreateTask)))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(`{
// 		"logins": [ "MyLogin1", "MyLogin2", "MyLogin3", "MyLogin4" ],
// 		"title": "MyTask1",
// 		"description": "This task is important one too!"
// 	}`)))

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	ma.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        initiatorLogin,
// 		}, nil)

// 	mts.On("InsertTask",
// 		&models.Task{
// 			ID:              taskId,
// 			Logins:          logins,
// 			ApprovalTokens:  approvalTokens,
// 			Title:           title,
// 			Description:     description,
// 			InitiatorLogin:  initiatorLogin,
// 			CurrApprovalNum: 0,
// 			Status:          models.TaskInProgressStatus,
// 		}).
// 		Return(fmt.Errorf("failed to insert task into DB"))

// 	createTaskHandler.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"error":"failed to create task requested by login test123"}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusInternalServerError, resp.StatusCode)
// }

// // GetTask

// func (s *unitTestSuite) TestGetTaskOk() {
// 	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
// 	logins := []string{
// 		"MyLogin1",
// 		"MyLogin2",
// 		"MyLogin3",
// 		"MyLogin4",
// 	}
// 	approvalTokens := []string{
// 		s.generateToken(taskId, "MyLogin1"),
// 		s.generateToken(taskId, "MyLogin2"),
// 		s.generateToken(taskId, "MyLogin3"),
// 		s.generateToken(taskId, "MyLogin4"),
// 	}
// 	title := "MyTask1"
// 	description := "This task is important one too!"
// 	initiatorLogin := "test123"

// 	logger, _ := zap.NewProduction()
// 	mts := new(mock_task_storage.MockTaskStorage)
// 	ma := new(mock_auth.MockAuth)
// 	m := mail.New(logger.Sugar())
// 	t := task.New(mts, m, ma, logger.Sugar(), func() uuid.UUID {
// 		taskIdUuid, _ := uuid.FromString(taskId)
// 		return taskIdUuid
// 	})
// 	srv := task_http.NewTest(t, logger.Sugar(), taskId)

// 	getTaskHandler := srv.AnnotateContext()(srv.ValidateAuth()(http.HandlerFunc(srv.GetTask)))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/tasks/"+taskId, nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	ma.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        initiatorLogin,
// 		}, nil)

// 	mts.On("GetTask", taskId).
// 		Return(&models.Task{
// 			ID:              taskId,
// 			Logins:          logins,
// 			ApprovalTokens:  approvalTokens,
// 			Title:           title,
// 			Description:     description,
// 			InitiatorLogin:  initiatorLogin,
// 			CurrApprovalNum: 0,
// 			Status:          models.TaskInProgressStatus,
// 		}, nil)

// 	getTaskHandler.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"ID":"` + taskId + `","Logins":["MyLogin1","MyLogin2","MyLogin3","MyLogin4"],"ApprovalTokens":["` + approvalTokens[0] + `","` + approvalTokens[1] + `","` + approvalTokens[2] + `","` + approvalTokens[3] + `"],"Title":"MyTask1","Description":"This task is important one too!","InitiatorLogin":"test123","CurrApprovalNum":0,"Status":0}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusOK, resp.StatusCode)
// }

// func (s *unitTestSuite) TestGetTaskNoTaskIdErr() {
// 	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
// 	initiatorLogin := "test123"

// 	logger, _ := zap.NewProduction()
// 	mts := new(mock_task_storage.MockTaskStorage)
// 	ma := new(mock_auth.MockAuth)
// 	m := mail.New(logger.Sugar())
// 	t := task.New(mts, m, ma, logger.Sugar(), func() uuid.UUID {
// 		taskIdUuid, _ := uuid.FromString(taskId)
// 		return taskIdUuid
// 	})
// 	srv := task_http.NewTest(t, logger.Sugar())

// 	getTaskHandler := srv.AnnotateContext()(srv.ValidateAuth()(http.HandlerFunc(srv.GetTask)))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/tasks/", nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	ma.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        initiatorLogin,
// 		}, nil)

// 	getTaskHandler.ServeHTTP(rec, req)

// 	resp := rec.Result()

// 	s.Equal(http.StatusNotFound, resp.StatusCode)
// }

// func (s *unitTestSuite) TestGetTaskGetTaskErrorErr() {
// 	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
// 	initiatorLogin := "test123"

// 	logger, _ := zap.NewProduction()
// 	mts := new(mock_task_storage.MockTaskStorage)
// 	ma := new(mock_auth.MockAuth)
// 	m := mail.New(logger.Sugar())
// 	t := task.New(mts, m, ma, logger.Sugar(), func() uuid.UUID {
// 		taskIdUuid, _ := uuid.FromString(taskId)
// 		return taskIdUuid
// 	})
// 	srv := task_http.NewTest(t, logger.Sugar(), taskId)

// 	getTaskHandler := srv.AnnotateContext()(srv.ValidateAuth()(http.HandlerFunc(srv.GetTask)))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/tasks/"+taskId, nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	ma.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        initiatorLogin,
// 		}, nil)

// 	mts.On("GetTask", taskId).
// 		Return(nil, fmt.Errorf("failed to get task from DB"))

// 	getTaskHandler.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"error":"failed to get info for task ID ` + taskId + `"}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusInternalServerError, resp.StatusCode)
// }

// func (s *unitTestSuite) TestGetTaskTaskNotFoundErr() {
// 	taskId := "1dea04bf-fd0c-48e0-9032-6ad3ddaea5af"
// 	initiatorLogin := "test123"

// 	logger, _ := zap.NewProduction()
// 	mts := new(mock_task_storage.MockTaskStorage)
// 	ma := new(mock_auth.MockAuth)
// 	m := mail.New(logger.Sugar())
// 	t := task.New(mts, m, ma, logger.Sugar(), func() uuid.UUID {
// 		taskIdUuid, _ := uuid.FromString(taskId)
// 		return taskIdUuid
// 	})
// 	srv := task_http.NewTest(t, logger.Sugar(), taskId)

// 	getTaskHandler := srv.AnnotateContext()(srv.ValidateAuth()(http.HandlerFunc(srv.GetTask)))

// 	rec := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/tasks/"+taskId, nil)

// 	req.AddCookie(&http.Cookie{
// 		Name:  "access",
// 		Value: testAccessToken,
// 	})
// 	req.AddCookie(&http.Cookie{
// 		Name:  "refresh",
// 		Value: testRefreshToken,
// 	})

// 	ma.On("ValidateAuth",
// 		&models.TokenPair{
// 			AccessToken:  testAccessToken,
// 			RefreshToken: testRefreshToken,
// 		}).
// 		Return(&models.AuthResult{
// 			Status:       models.OkAuthRespStatus,
// 			AccessToken:  "",
// 			RefreshToken: "",
// 			Login:        initiatorLogin,
// 		}, nil)

// 	mts.On("GetTask", taskId).
// 		Return(&models.Task{}, nil)

// 	getTaskHandler.ServeHTTP(rec, req)

// 	resp := rec.Result()
// 	defer resp.Body.Close()

// 	data, err := ioutil.ReadAll(resp.Body)
// 	s.Nil(err, "cannot read data from response")

// 	expected := `{"error":"no task found with task ID ` + taskId + `"}`

// 	s.Equal(expected, string(data))
// 	s.Equal(http.StatusNotFound, resp.StatusCode)
// }
