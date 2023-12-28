package resolvers

import (
	"errors"
	"log"

	"github.com/muety/telepush/api"
	"github.com/muety/telepush/model"
)

type TextResolver struct{}

func (r TextResolver) IsValid(m *model.Message) error {
	if len(m.Text) == 0 {
		return errors.New("text parameter missing")
	}
	return nil
}

func (r TextResolver) Resolve(recipientId string, m *model.Message, params *model.MessageOptions) error {
	defer logMessage(m)
	var disableLinkPreview bool
	if params != nil {
		disableLinkPreview = params.DisableLinkPreviews
	}

	err := api.SendMessage(&model.TelegramOutMessage{
		ChatId:             recipientId,
		Text:               m.Text,
		ParseMode:          params.ParseMode(),
		DisableLinkPreview: disableLinkPreview,
	})

	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return err
}
