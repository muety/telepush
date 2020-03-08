package internal

import (
	"errors"
	"github.com/n1try/telegram-middleman-bot/api"
	"github.com/n1try/telegram-middleman-bot/model"
)

func validateText(m *model.InMessage) error {
	if len(m.Text) == 0 {
		return errors.New("text parameter missing")
	}
	return nil
}

func logText(m *model.InMessage) string {
	return m.Text
}

func resolveText(recipientId string, m *model.InMessage) error {
	return api.SendMessage(recipientId, "*"+m.Origin+"* wrote:\n\n"+m.Text)
}
