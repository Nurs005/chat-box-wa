package handler

import (
	"encoding/json"
	"errors"
	"github.com/chatbox/whatsapp/internal/domain"
	errors2 "github.com/chatbox/whatsapp/pkg/errors"
	"github.com/chatbox/whatsapp/pkg/utils"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(SessionCtxKey).(*domain.Session)

	if !ok || session == nil {
		utils.WriteJSONResponse(w, r, http.StatusUnauthorized, "session not found", false)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("WebSocket upgrade failed")
		utils.WriteJSONResponse(w, r, http.StatusBadRequest, "smth wrong with upgrader", false)
		return
	}

	log.Info().Str("token", session.Token).Msg("WebSocket client connected")

	h.service.IHub.HandleConnection(session.Token, conn)

	defer func() {
		h.service.DeleteConnection(session.Token)
		log.Info().Str("token", session.Token).Msg("WebSocket client disconnected")
	}()

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("WebSocket read error")
			break
		}
		log.Info().Str("token", session.Token).Str("msg", string(msg)).Msg("Received WebSocket message")

		switch t {
		case websocket.TextMessage:
			if string(msg) == "ping" {
				h.service.IHub.Send(session.Token, []byte("pong"), websocket.TextMessage)
			} else {
				data, err := json.Marshal(&domain.BaseResponse{Success: false, StatusCode: http.StatusBadRequest, Message: errors2.InvalidTypeOfTextMsg.Error()})
				if err != nil {
					log.Error().Err(err).Msg("WebSocket marshal error")
					return
				}
				h.service.IHub.Send(session.Token, data, websocket.BinaryMessage)
			}
		case websocket.BinaryMessage:
			err := h.service.IWebSocketCommandResolver.Handle(session, msg)
			if err != nil {
				log.Error().Err(err).Msg("WebSocket command resolver error")

				var statusCode int

				switch {
				case errors.Is(err, errors2.InvalidCommand):
					statusCode = http.StatusBadRequest

				case errors.Is(err, errors2.InvalidBodyForCmd):
					statusCode = http.StatusBadRequest
				case strings.Contains(err.Error(), "invalid JID:"):
					statusCode = http.StatusBadRequest
				default:
					statusCode = http.StatusInternalServerError
				}

				data, err := json.Marshal(&domain.BaseResponse{Success: false, StatusCode: statusCode, Message: err.Error()})
				if err != nil {
					log.Error().Err(err).Msg("WebSocket marshal error")
					return
				}

				h.service.IHub.Send(session.Token, data, websocket.BinaryMessage)
			}
		}

	}
}
