package config

import (
	"github.com/leandro-lugaresi/hub"
)

var eventHub *hub.Hub

func GetHub() *hub.Hub {
	if eventHub == nil {
		eventHub = hub.New()
	}
	return eventHub
}
