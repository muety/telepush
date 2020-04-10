package inlets

import "net/http"

type Inlet interface {
	Handler(http.Handler) http.Handler
}
