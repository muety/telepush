package inlets

import (
	"context"
	"encoding/json"
	"github.com/n1try/telegram-middleman-bot/config"
	"github.com/n1try/telegram-middleman-bot/model"
	"net/http"
)

type DefaultInlet struct{}

func NewDefaultInlet() *DefaultInlet {
	return &DefaultInlet{}
}

func (d DefaultInlet) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var m model.InMessage
		err := dec.Decode(&m)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		next(
			w,
			r.WithContext(context.WithValue(r.Context(), config.KeyMessage, &m)),
		)
	}
}
