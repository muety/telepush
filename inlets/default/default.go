package _default

import (
	"context"
	"encoding/json"
	"github.com/muety/webhook2telegram/config"
	"github.com/muety/webhook2telegram/model"
	"github.com/muety/webhook2telegram/util"
	"net/http"

	"github.com/muety/webhook2telegram/inlets"
)

type DefaultInlet struct{}

func (i *DefaultInlet) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m model.ExtendedMessage

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if len(m.Origin) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing origin parameter"))
			return
		}

		m.Text = "*" + util.EscapeMarkdown(m.Origin) + "* wrote:\n\n" + m.Text

		options := &model.MessageParams{}
		if m.Options != nil {
			options = m.Options
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, config.KeyMessage, &m.DefaultMessage)
		ctx = context.WithValue(ctx, config.KeyParams, options)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func New() inlets.Inlet {
	return &DefaultInlet{}
}
