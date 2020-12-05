package middleware

import (
	"github.com/leandro-lugaresi/hub"
	"github.com/muety/webhook2telegram/config"
	"net/http"
)

func WithEventLogging() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &StatusRecorderWriter{
				ResponseWriter: w,
			}

			h.ServeHTTP(recorder, r)

			if recorder.IsSuccess() {
				config.GetHub().Publish(hub.Message{
					Name: config.EventOnRequestSuccessful,
				})
			} else {
				config.GetHub().Publish(hub.Message{
					Name: config.EventOnRequestFailed,
				})
			}
		})
	}
}
