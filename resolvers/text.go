package resolvers

import (
	"errors"
	"github.com/n1try/telegram-middleman-bot/api"
	"github.com/n1try/telegram-middleman-bot/model"
)

func validateText(m *model.DefaultMessage) error {
	if len(m.Text) == 0 {
		return errors.New("text parameter missing")
	}
	return nil
}

func logText(m *model.DefaultMessage) string {
	return m.Text
}

func resolveText(recipientId string, m *model.DefaultMessage, params *model.MessageParams) error {
	var disableLinkPreview bool
	if params != nil {
		disableLinkPreview = params.DisableLinkPreviews
	}

	return api.SendMessage(&model.TelegramOutMessage{
		ChatId:             recipientId,
		Text:               "*" + m.Origin + "* wrote:\n\n" + m.Text,
		ParseMode:          "Markdown",
		DisableLinkPreview: disableLinkPreview,
	})
}
