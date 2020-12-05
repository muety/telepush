package model

type StoreObject struct {
	User   TelegramUser `json:"user"`
	ChatId int          `json:"chat_id"`
}
