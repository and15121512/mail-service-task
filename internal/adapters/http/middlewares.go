package http

import (
	"context"
	"errors"
	"fmt"
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

			accessTokenCookie, err := r.Cookie("access")
			if err != nil {
				utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
					"error": tokenReadingFailed,
				})
				logger.Errorf(tokenReadingFailed)
				return
			}
			refreshTokenCookie, err := r.Cookie("refresh")
			refreshToken := refreshTokenCookie.Value
			if errors.Is(err, http.ErrNoCookie) {
				refreshToken = ""
			} else if err != nil {
				utils.ResponseJSON(w, http.StatusForbidden, map[string]string{
					"error": tokenReadingFailed,
				})
				logger.Errorf(tokenReadingFailed)
				return
			}

			ar, err := s.task.ValidateAuth(r.Context(), &models.TokenPair{
				AccessToken:  accessTokenCookie.Value,
				RefreshToken: refreshToken,
			})
			if err != nil {
				utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
					"error": fmt.Sprintf("failed to validate token with auth service: %s", err.Error()),
				})
				logger.Errorf("failed to validate token with auth service: %s", err.Error())
				return
			}
			if ar.Status == models.RefusedAuthRespStatus {
				utils.ResponseJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "authorization required",
				})
				logger.Errorf("authorization required")
				return
			} else if ar.Status == models.RefreshedAuthRespStatus {
				http.SetCookie(w, &http.Cookie{
					Name:  "access",
					Value: ar.AccessToken,
				})
				http.SetCookie(w, &http.Cookie{
					Name:  "refresh",
					Value: ar.RefreshToken,
				})
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxKeyUser{}, ar.Login)
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
