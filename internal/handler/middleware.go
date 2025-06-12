package handler

import (
	"context"
	"github.com/chatbox/whatsapp/pkg/utils"
	"net/http"
	"strings"
)

type ctxKey string

const SessionCtxKey ctxKey = "session"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteJSONResponse(w, r, http.StatusUnauthorized, "unauthorized")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		session, err := h.service.ISessionService.GetSession(r.Context(), token)
		if err != nil {
			utils.WriteJSONResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), SessionCtxKey, session)

		newR := r.Clone(ctx)

		next.ServeHTTP(w, newR)
	})
}
