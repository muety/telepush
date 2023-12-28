package handlers

import (
	"errors"
	"github.com/muety/telepush/config"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/resolvers"
	"github.com/muety/telepush/services"
	"github.com/muety/telepush/util"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
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
	var m *model.Message
	var p model.MessageOptions

	if message := r.Context().Value(config.KeyMessage); message != nil {
		m = message.(*model.Message)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse message"))
		return
	}

	if utf8.RuneCountInString(m.Text) > 4096 {
		if !config.Get().TruncateMsgs {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("message too long (max is 4096 characters)"))
			return
		}
		m.Text = util.TruncateInRunes(m.Text, 4096)
	}

	if params := r.Context().Value(config.KeyParams); params != nil {
		p = *(params.(*model.MessageOptions))
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
		statusCode := parseStatusCode(err, http.StatusInternalServerError)

		if statusCode == http.StatusForbidden {
			// user has probably blocked the bot -> invalidate token
			h.userService.InvalidateToken(token)
			log.Printf("invalidating token '%s' for chat '%s', because got 403 from telegram\n", token, recipientId)
			err = errors.New("error: got 403 from telegram, invalidating your token, text the bot to generate a new one")
		}

		w.WriteHeader(statusCode)
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
