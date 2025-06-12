package handler

import (
	"errors"
	"github.com/chatbox/whatsapp/internal/domain"
	er "github.com/chatbox/whatsapp/pkg/errors"
	"github.com/chatbox/whatsapp/pkg/utils"
	"github.com/go-chi/render"
	"net/http"
)

func (h *Handler) CreateSessionForUser(w http.ResponseWriter, r *http.Request) {
	session := new(domain.Session)
	if err := render.DecodeJSON(r.Body, session); err != nil {
		utils.WriteJSONResponse(w, r, http.StatusBadRequest, err.Error(), false)
		return
	}

	if err := h.service.ISessionService.NewSession(r.Context(), session); err != nil {
		if errors.Is(err, er.HaveNoBusinessId) {
			utils.WriteJSONResponse(w, r, http.StatusBadRequest, err.Error(), false)
		}
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusCreated, session, true)
}

func (h *Handler) GetQR(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(SessionCtxKey).(*domain.Session)

	if !ok || session == nil {
		utils.WriteJSONResponse(w, r, http.StatusUnauthorized, "session not found", false)
		return
	}

	code, err := h.service.ISessionService.GetQRLogin(r.Context(), session.Token)
	if err != nil {
		utils.WriteJSONResponse(w, r, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, code, true)
}
