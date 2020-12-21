package resolvers

import (
	"errors"
	"log"

	"github.com/muety/webhook2telegram/api"
	"github.com/muety/webhook2telegram/model"
)

type TextResolver struct{}

func (r TextResolver) IsValid(m *model.DefaultMessage) error {
	if len(m.Text) == 0 {
		return errors.New("text parameter missing")
	}
	return nil
}

func (r TextResolver) Resolve(recipientId string, m *model.DefaultMessage, params *model.MessageParams) error {
	defer logMessage(m)
	var disableLinkPreview bool
	if params != nil {
		disableLinkPreview = params.DisableLinkPreviews
	}

	err := api.SendMessage(&model.TelegramOutMessage{
		ChatId:             recipientId,
		Text:               m.Text,
		ParseMode:          "Markdown",
		DisableLinkPreview: disableLinkPreview,
	})

	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return err
}
