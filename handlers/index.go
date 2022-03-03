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
	return &IndexHandler{
		Tpl:    template.Must(template.ParseFiles("views/index.tpl.html")),
		Config: config.Get(),
	}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Tpl.Execute(w, indexData{Config: h.Config})
}
