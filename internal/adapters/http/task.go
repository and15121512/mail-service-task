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
	r.With(s.AnnotateContext()).With(s.ValidateAuth()).Get("/tasks/{task_id}", s.GetTask)
	r.With(s.AnnotateContext()).With(s.ValidateAuth()).Post("/tasks/{task_id}", s.UpdateTask)
	r.With(s.AnnotateContext()).With(s.ValidateAuth()).Delete("/tasks/{task_id}", s.DeleteTask)

	r.With(s.AnnotateContext()).With(s.ValidateAuth()).Post("/tasks/{task_id}/approve", s.ApproveOrDecline)

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

	login, ok := r.Context().Value(ctxKeyUser{}).(string)
	if !ok {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to extract token from request",
		})
		logger.Errorf("failed to extract token from request")
		return
	}
	user := &models.User{Login: login}

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

func (s *Server) GetTask(w http.ResponseWriter, r *http.Request) {
	logger := s.annotatedLogger(r.Context())

	task_id := chi.URLParam(r, "task_id")
	if task_id == "" {
		utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
			"error": "No task ID provided in URL for get task request",
		})
		logger.Errorf("No task ID provided in URL for get task request")
		return
	}

	task, err := s.task.GetTask(r.Context(), task_id)
	if err != nil {
		utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
			"error": fmt.Sprintf("failed to get info for task ID %s", task_id),
		})
		logger.Errorf("failed to get info for task ID %s", task_id)
		return
	}
	if task.ID == "" {
		utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
			"error": fmt.Sprintf("no task found with task ID %s", task_id),
		})
		logger.Infof("no task found with task ID %s", task_id)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to encode task requested with task ID %s", task_id),
		})
		logger.Errorf("failed to encode task requested with task ID %s", task_id)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateTaskRequest struct {
	Logins      []string `json:"logins"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
}

func (s *Server) UpdateTask(w http.ResponseWriter, r *http.Request) {
	logger := s.annotatedLogger(r.Context())
	task_id := chi.URLParam(r, "task_id")
	if task_id == "" {
		utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
			"error": "No task ID provided in URL for update task request",
		})
		logger.Errorf("No task ID provided in URL for update task request")
		return
	}
	var req updateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cannot parse input for update task request",
		})
		logger.Errorf("cannot parse input for update task request: %s", err.Error())
		return
	}

	login, ok := r.Context().Value(ctxKeyUser{}).(string)
	if !ok {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to extract token from request",
		})
		logger.Errorf("failed to extract token from request")
		return
	}
	user := &models.User{Login: login}

	err = s.task.UpdateTask(r.Context(), models.Task{
		ID:          task_id,
		Logins:      req.Logins,
		Title:       req.Title,
		Description: req.Description,
	}, *user)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("failed to update task requested by login %s; are you the task author?", user.Login),
		})
		logger.Errorf("failed to update task requested by login %s", user.Login)
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]string{
		"task_id": task_id,
	})
}

func (s *Server) DeleteTask(w http.ResponseWriter, r *http.Request) {
	logger := s.annotatedLogger(r.Context())

	task_id := chi.URLParam(r, "task_id")
	if task_id == "" {
		utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
			"error": "No task ID provided in URL for delete task request",
		})
		logger.Errorf("No task ID provided in URL for delete task request")
		return
	}

	login, ok := r.Context().Value(ctxKeyUser{}).(string)
	if !ok {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to extract token from request",
		})
		logger.Errorf("failed to extract token from request")
		return
	}
	user := &models.User{Login: login}

	err := s.task.DeleteTask(r.Context(), task_id, *user)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("failed to delete task requested by login %s; are you the task author?", user.Login),
		})
		logger.Errorf("failed to delete task requested by login %s", user.Login)
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]string{})
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) ApproveOrDecline(w http.ResponseWriter, r *http.Request) {
	logger := s.annotatedLogger(r.Context())

	task_id := chi.URLParam(r, "task_id")
	if task_id == "" {
		utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
			"error": "no task ID provided in URL for approve request",
		})
		logger.Errorf("no task ID provided in URL for approve request")
		return
	}
	token := r.URL.Query().Get("token")
	decision := r.URL.Query().Get("decision")
	if token == "" || decision == "" {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": "no token or no decision provided in URL query params for approve request",
		})
		logger.Errorf("no token or no decision provided in URL query params for approve request")
		return
	}

	err := s.task.ApproveOrDecline(r.Context(), task_id, token, decision)
	if err != nil {
		utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
			"error": "failed to process task with given task ID, token and decision",
		})
		logger.Errorf("failed to process task with given task ID, token and decision")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]string{})
}
