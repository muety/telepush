package config

import (
	"github.com/muety/webhook2telegram/store"
)

var storeInstance store.Store

func GetStore() store.Store {
	if storeInstance == nil {
		storeInstance = store.NewGobStore(Get().GetStorePath())
	}
	return storeInstance
}
