package inlets

import "net/http"

type Inlet interface {
	Middleware(http.HandlerFunc) http.HandlerFunc
}
