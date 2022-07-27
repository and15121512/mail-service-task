package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
)

const (
	tokenReadingFailed = "cannot read token from header"
	invalidToken       = "invalid token"
)

type ctxKeyUser struct{}

func (s *Server) ValidateAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := s.annotatedLogger(r.Context())

			accessToken, err := r.Cookie("access")
			if err != nil {
				utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
					"error": tokenReadingFailed,
				})
				logger.Errorf(tokenReadingFailed)
				return
			}

			// GRPC call to auth mocked!
			user, err := &models.User{
				Login:        "test123",
				PasswordHash: "b1b3773a05c0ed0176787a4f1574ff0075f7521e",
			}, nil
			_ = accessToken.Value
			// GRPC call to auth mocked!

			if err != nil {
				utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
					"error": invalidToken,
				})
				logger.Errorf(invalidToken)
				return
			}
			ctx := r.Context()

			ctx = context.WithValue(ctx, ctxKeyUser{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *Server) AnnotateContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, utils.CtxKeyRequestIDGet(), middleware.GetReqID(r.Context()))
			ctx = context.WithValue(ctx, utils.CtxKeyMethodGet(), r.Method)
			ctx = context.WithValue(ctx, utils.CtxKeyURLGet(), r.URL.String())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
