package handlers

import (
	"github.com/muety/telepush/config"
	"html/template"
	"net/http"
)

type IndexHandler struct {
	Tpl    *template.Template
	Config *config.BotConfig
}
type indexData struct {
	Config *config.BotConfig
}

func NewIndexHandler() *IndexHandler {
	h := &IndexHandler{Config: config.Get()}
	h.loadTemplates()
	return h
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Config.Env == "dev" {
		h.loadTemplates()
	}
	h.Tpl.Execute(w, indexData{Config: h.Config})
}

func (h *IndexHandler) loadTemplates() {
	h.Tpl = template.Must(template.ParseFiles("views/index.tpl.html"))
}
