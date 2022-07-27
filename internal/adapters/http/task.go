package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
)

func (s *Server) taskHandlers() http.Handler {
	r := chi.NewRouter()
	r.With(s.AnnotateContext()).With(s.ValidateAuth()).Post("/tasks", s.CreateTask)

	// Here are some other endpoints...

	return r
}

func (s *Server) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return s.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

type createTaskRequest struct {
	Logins      []string `json:"logins"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
}

func (s *Server) CreateTask(w http.ResponseWriter, r *http.Request) {
	logger := s.annotatedLogger(r.Context())
	var req createTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cannot parse input for create task request",
		})
		logger.Errorf("cannot parse input for create task request: %s", err.Error())
		return
	}

	user, ok := r.Context().Value(ctxKeyUser{}).(*models.User)
	if !ok {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to extract token from request",
		})
		logger.Errorf("failed to extract token from request")
		return
	}

	task, err := s.task.CreateTask(r.Context(), models.Task{
		Logins:         req.Logins,
		Title:          req.Title,
		Description:    req.Description,
		InitiatorLogin: user.Login,
	})
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to create task requested by login %s", user.Login),
		})
		logger.Errorf("failed to create task requested by login %s", user.Login)
		return
	}
	w.Header().Set("Location", r.URL.String()+"/"+task.ID)
	utils.ResponseJSON(w, http.StatusCreated, map[string]string{
		"task_id": task.ID,
	})
}

func GetTask(w http.ResponseWriter, r *http.Request) {
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
}

func ApproveOrDecline(w http.ResponseWriter, r *http.Request) {
}
