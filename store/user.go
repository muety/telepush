package store

import (
	"github.com/n1try/telegram-middleman-bot/model"
	"strconv"
)

func InvalidateToken(userChatId int) {
	for k, v := range GetMap() {
		entry, ok := v.(model.StoreObject)
		if ok && entry.ChatId == userChatId {
			Delete(k)
		}
	}
}

func ResolveToken(token string) string {
	value := Get(token)
	if value != nil {
		return strconv.Itoa((value.(model.StoreObject)).ChatId)
	}
	return ""
}
