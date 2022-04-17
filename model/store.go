package model

type StoreObject struct {
	User   TelegramUser `json:"user"`
	ChatId int64        `json:"chat_id"`
}
