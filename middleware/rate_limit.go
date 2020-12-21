package middleware

import (
	"fmt"
	"github.com/muety/webhook2telegram/config"
	"github.com/n1try/limiter/v3"
	mhttp "github.com/n1try/limiter/v3/drivers/middleware/stdlib"
	memst "github.com/n1try/limiter/v3/drivers/store/memory"
	"net/http"
)

func WithRateLimit() func(h http.Handler) http.Handler {
	rate, _ := limiter.NewRateFromFormatted(fmt.Sprintf("%d-H", config.Get().ReqRateLimit))
	store := memst.NewStore()
	l := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	return mhttp.NewMiddleware(l).Handler
}
