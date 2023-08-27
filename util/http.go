package util

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

type RouterSwapper struct {
	Root   *mux.Router
	Prefix string
	mu     sync.Mutex
}

func (rs *RouterSwapper) Swap(newRouter *mux.Router) {
	rs.mu.Lock()
	rs.Root = newRouter
	rs.mu.Unlock()
}

func (rs *RouterSwapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rs.mu.Lock()
	root := rs.Root
	rs.mu.Unlock()
	root.ServeHTTP(w, r)
}
