package services

import (
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/store"
	"strconv"
)

type UserService struct {
	store store.Store
}

func NewUserService(store store.Store) *UserService {
	return &UserService{store: store}
}

func (s *UserService) InvalidateToken(userChatId int) {
	for k, v := range s.store.GetItems() {
		entry, ok := v.(model.StoreObject)
		if ok && entry.ChatId == userChatId {
			s.store.Delete(k)
		}
	}
}

func (s *UserService) ResolveToken(token string) string {
	value := s.store.Get(token)
	if value != nil {
		return strconv.Itoa((value.(model.StoreObject)).ChatId)
	}
	return ""
}
