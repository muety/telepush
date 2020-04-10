package middleware

import (
	"fmt"
	limiter "github.com/n1try/limiter/v3"
	mhttp "github.com/n1try/limiter/v3/drivers/middleware/stdlib"
	memst "github.com/n1try/limiter/v3/drivers/store/memory"
	"github.com/n1try/telegram-middleman-bot/config"
	"net/http"
)

func NewCheckMethod(cfg *config.BotConfig) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				w.WriteHeader(415)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func NewRateLimit(cfg *config.BotConfig) func(h http.Handler) http.Handler {
	rate, _ := limiter.NewRateFromFormatted(fmt.Sprintf("%d-H", cfg.RateLimit))
	store := memst.NewStore()
	l := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	return mhttp.NewMiddleware(l).Handler
}
