package _default

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/n1try/telegram-middleman-bot/config"
	"github.com/n1try/telegram-middleman-bot/inlets"
	"github.com/n1try/telegram-middleman-bot/model"
)

type DefaultInlet struct {
	inlets.Inlet
}

func New() inlets.Inlet {
	return &DefaultInlet{}
}

func (i *DefaultInlet) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m model.DefaultMessage

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&m); err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		if len(m.Origin) == 0 {
			w.WriteHeader(400)
			w.Write([]byte("missing origin parameter"))
			return
		}

		m.Text = "*" + m.Origin + "* wrote:\n\n" + m.Text

		next(
			w,
			r.WithContext(context.WithValue(r.Context(), config.KeyMessage, &m)),
		)
	}
}
