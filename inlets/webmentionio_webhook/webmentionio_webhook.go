package webmentionio_webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/muety/telepush/config"
	"github.com/muety/telepush/inlets"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/resolvers"
	"github.com/muety/telepush/util"
	"net/http"
	"net/url"
)

type WebmentionioInlet struct{}

func New() inlets.Inlet {
	return &WebmentionioInlet{}
}

func (i *WebmentionioInlet) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m WebmentionMessage

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if !validateMessage(&m) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		message := transformMessage(&m)

		ctx := r.Context()
		ctx = context.WithValue(ctx, config.KeyMessage, message)
		ctx = context.WithValue(ctx, config.KeyParams, &model.MessageParams{DisableLinkPreviews: true})

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateMessage(message *WebmentionMessage) bool {
	if u, err := url.Parse(message.Source); err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}
	if u, err := url.Parse(message.Target); err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}

	return message.Secret != ""
}

func transformMessage(in *WebmentionMessage) *model.DefaultMessage {
	text := "*Webmention Watcher* wrote:\n\n"
	text += util.EscapeMarkdown(fmt.Sprintf("Your article at %s was mentioned at %s.", in.Target, in.Source))

	return &model.DefaultMessage{
		RecipientToken: in.Secret,
		Text:           text,
		Type:           resolvers.TextType,
	}
}
