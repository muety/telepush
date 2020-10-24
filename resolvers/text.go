package resolvers

import (
	"errors"
	"log"

	"github.com/n1try/telegram-middleman-bot/api"
	"github.com/n1try/telegram-middleman-bot/model"
)

type TextResolver struct{}

func (r TextResolver) IsValid(m *model.DefaultMessage) error {
	if len(m.Text) == 0 {
		return errors.New("text parameter missing")
	}
	return nil
}

func (r TextResolver) Resolve(recipientId string, m *model.DefaultMessage, params *model.MessageParams) *model.ApiError {
	var disableLinkPreview bool
	if params != nil {
		disableLinkPreview = params.DisableLinkPreviews
	}

	apiErr := api.SendMessage(&model.TelegramOutMessage{
		ChatId:             recipientId,
		Text:               m.Text,
		ParseMode:          "Markdown",
		DisableLinkPreview: disableLinkPreview,
	})

	if apiErr != nil {
		log.Printf("error: %v\n", apiErr)
	}

	return apiErr
}
