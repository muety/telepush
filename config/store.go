package config

import (
	"github.com/muety/telepush/store"
)

var storeInstance store.Store

func GetStore() store.Store {
	if storeInstance == nil {
		storeInstance = store.NewGobStore(Get().GetStorePath())
	}
	return storeInstance
}
