package resolvers

import (
	b64 "encoding/base64"
	"errors"
	"github.com/n1try/telegram-middleman-bot/api"
	"github.com/n1try/telegram-middleman-bot/model"
	"net/http"
)

func validateFile(m *model.DefaultMessage) error {
	if len(m.File) == 0 || len(m.Filename) == 0 {
		return errors.New("file or file name parameter missing")
	}
	return nil
}

func logFile(m *model.DefaultMessage) string {
	return "A document named " + m.Filename + " was sent"
}

func resolveFile(recipientId string, m *model.DefaultMessage, params *model.MessageParams) *model.ApiError {
	decodedFile, err := b64.StdEncoding.DecodeString(m.File)
	if err != nil {
		return &model.ApiError{
			StatusCode: http.StatusBadRequest,
			Text:       err.Error(),
		}
	}

	return api.SendDocument(&model.TelegramOutDocument{
		ChatId:    recipientId,
		Caption:   "*" + m.Origin + "* sent a document",
		ParseMode: "Markdown",
		Document: &model.TelegramInputFile{
			Name: m.Filename,
			Data: decodedFile,
		},
	})
}
