package handlers

import (
	"github.com/muety/webhook2telegram/config"
	"github.com/muety/webhook2telegram/model"
	"github.com/muety/webhook2telegram/resolvers"
	"github.com/muety/webhook2telegram/services"
	"net/http"
)

type MessageHandler struct {
	userService *services.UserService
}

func NewMessageHandler(userService *services.UserService) *MessageHandler {
	return &MessageHandler{userService: userService}
}

func (h *MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m *model.DefaultMessage
	var p model.MessageParams

	if message := r.Context().Value(config.KeyMessage); message != nil {
		m = message.(*model.DefaultMessage)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse message"))
		return
	}

	if params := r.Context().Value(config.KeyParams); params != nil {
		p = params.(model.MessageParams)
	}

	token := r.Header.Get("token")
	if token == "" {
		token = m.RecipientToken
	}

	if len(token) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing recipient_token parameter"))
		return
	}

	// TODO: Refactoring: get rid of this resolver concept
	resolver := resolvers.GetResolver(m.Type)

	if err := resolver.IsValid(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	recipientId := h.userService.ResolveToken(token)

	if len(recipientId) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("passed token does not relate to a valid user"))
		return
	}

	go resolver.Resolve(recipientId, m, &p)

	w.WriteHeader(http.StatusAccepted)
}
