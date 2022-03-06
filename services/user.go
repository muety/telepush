package services

import (
	"fmt"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/store"
	"strconv"
)

type UserService struct {
	store store.Store
}

type Tokens []string

func (tokens Tokens) String() string {
	var str string
	for i, t := range tokens {
		str += fmt.Sprintf("*%d:* `%s`\n", i+1, t)
	}
	return str
}

func NewUserService(store store.Store) *UserService {
	return &UserService{store: store}
}

func (s *UserService) SetToken(token string, fromUser model.TelegramUser, chatId int) {
	s.store.Put(token, model.StoreObject{User: fromUser, ChatId: chatId})
}

func (s *UserService) InvalidateToken(token string) {
	s.store.Delete(token)
}

func (s *UserService) ResolveToken(token string) string {
	value := s.store.Get(token)
	if value != nil {
		return strconv.Itoa((value.(model.StoreObject)).ChatId)
	}
	return ""
}

// O(n)
func (s *UserService) ListTokens(chatId int) Tokens {
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
