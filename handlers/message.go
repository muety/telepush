package handlers

import (
	"github.com/muety/telepush/config"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/resolvers"
	"github.com/muety/telepush/services"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var telegramApiErrorRegexp *regexp.Regexp

func init() {
	telegramApiErrorRegexp = regexp.MustCompile(`telegram api returned status (\d{3}):`)
}

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
		p = *(params.(*model.MessageParams))
	}

	var token string
	if t := r.Context().Value(config.KeyRecipient); t != nil {
		token = t.(string)
	}

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing recipient token"))
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

	var async bool
	if asyncParam := r.URL.Query().Get("async"); strings.ToLower(asyncParam) == "true" || asyncParam == "1" {
		async = true
	}

	if async {
		go resolver.Resolve(recipientId, m, &p)
		w.WriteHeader(http.StatusAccepted)
		return
	} else if err := resolver.Resolve(recipientId, m, &p); err != nil {
		w.WriteHeader(parseStatusCode(err, http.StatusInternalServerError))
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// very hacky!
func parseStatusCode(err error, fallback int) int {
	if matches := telegramApiErrorRegexp.FindStringSubmatch(err.Error()); len(matches) > 1 {
		if statusCode, err := strconv.Atoi(matches[1]); err == nil {
			return statusCode
		}
		return fallback
	}
	return fallback
}
