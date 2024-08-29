package services

import (
	"fmt"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/store"
	"maps"
	"slices"
	"strconv"
	"sync"
)

type UserService struct {
	store        store.Store
	usersByChats map[int64][]int64
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
		instance = &UserService{store: store, usersByChats: make(map[int64][]int64)}
	})
	return instance
}

func (s *UserService) SetToken(token string, fromUser model.TelegramUser, chatId int64) {
	clear(s.usersByChats)
	s.store.Put(token, model.StoreObject{User: fromUser, ChatId: chatId})
}

func (s *UserService) InvalidateToken(token string) {
	clear(s.usersByChats)
	s.store.Delete(token)
}

func (s *UserService) ResolveToken(token string) string {
	value := s.store.Get(token)
	if value != nil {
		return strconv.FormatInt((value.(model.StoreObject)).ChatId, 10)
	}
	return ""
}

func (s *UserService) GetUsers() []int64 {
	usersMap := map[int64]bool{}
	for _, v := range s.store.GetItems() {
		if obj, ok := v.(model.StoreObject); ok {
			usersMap[obj.User.Id] = true
		}
	}
	return slices.Collect(maps.Keys(usersMap))
}

func (s *UserService) GetChats() []int64 {
	chatsMap := map[int64]bool{}
	for _, v := range s.store.GetItems() {
		if obj, ok := v.(model.StoreObject); ok {
			chatsMap[obj.ChatId] = true
		}
	}
	return slices.Collect(maps.Keys(chatsMap))
}

func (s *UserService) GetChatsStr() []string {
	chats := s.GetChats()
	chatsStr := make([]string, len(chats))
	for i, c := range chats {
		chatsStr[i] = strconv.FormatInt(c, 10)
	}
	return chatsStr
}

func (s *UserService) GetUsersByChat(chatId int64) []int64 {
	if users, ok := s.usersByChats[chatId]; ok {
		return users
	}

	userIds := make([]int64, 0)
	for _, v := range s.store.GetItems() {
		if obj, ok := v.(model.StoreObject); ok {
			if obj.ChatId == chatId {
				userIds = append(userIds, obj.User.Id)
			}
		}
	}

	s.usersByChats[chatId] = userIds
	return userIds
}

func (s *UserService) GetUsersByRecipient(recipientId string) []int64 {
	chatId, err := strconv.ParseInt(recipientId, 10, 64)
	if err != nil {
		return []int64{}
	}
	return s.GetUsersByChat(chatId)
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
