package services

import (
	"fmt"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/store"
	"strconv"
	"sync"
)

type UserService struct {
	store store.Store
}

var instance *UserService // singleton
var once sync.Once

type Tokens []string

func (tokens Tokens) String() string {
	var str string
	for i, t := range tokens {
		str += fmt.Sprintf("*%d:* `%s`\n", i+1, t)
	}
	return str
}

func NewUserService(store store.Store) *UserService {
	once.Do(func() {
		instance = &UserService{store: store}
	})
	return instance
}

func (s *UserService) SetToken(token string, fromUser model.TelegramUser, chatId int64) {
	s.store.Put(token, model.StoreObject{User: fromUser, ChatId: chatId})
}

func (s *UserService) InvalidateToken(token string) {
	s.store.Delete(token)
}

func (s *UserService) ResolveToken(token string) string {
	value := s.store.Get(token)
	if value != nil {
		return strconv.FormatInt((value.(model.StoreObject)).ChatId, 10)
	}
	return ""
}

// O(n)
func (s *UserService) ListTokens(chatId int64) Tokens {
	tokens := make(Tokens, 0)
	for k, v := range s.store.GetItems() {
		if obj, ok := v.(model.StoreObject); ok {
			if obj.ChatId == chatId {
				tokens = append(tokens, k)
			}
		}
	}
	return tokens
}
