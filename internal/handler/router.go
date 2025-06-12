package handler

import (
	"github.com/chatbox/whatsapp/internal"
	"github.com/chatbox/whatsapp/internal/infrastructure"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	log     *infrastructure.Logger
	service *internal.WuzApi
}

func NewHandler(log *infrastructure.Logger, service *internal.WuzApi) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // для тестов — потом обезопасим
	},
}

func (h *Handler) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/session/create", h.CreateSessionForUser)

	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		r.Post("/ws", h.HandleWebSocket)

		r.Get("/session/qr", h.GetQR)
	})

	return r
}
