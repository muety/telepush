package resolvers

import (
	b64 "encoding/base64"
	"errors"
	"github.com/n1try/telegram-middleman-bot/api"
	"github.com/n1try/telegram-middleman-bot/model"
	"log"
	"net/http"
)

type FileResolver struct{}

func (r FileResolver) IsValid(m *model.DefaultMessage) error {
	if len(m.File) == 0 || len(m.Filename) == 0 {
		return errors.New("file or file name parameter missing")
	}
	return nil
}

func (r FileResolver) Resolve(recipientId string, m *model.DefaultMessage, params *model.MessageParams) *model.ApiError {
	decodedFile, err := b64.StdEncoding.DecodeString(m.File)
	if err != nil {
		return &model.ApiError{
			StatusCode: http.StatusBadRequest,
			Text:       err.Error(),
		}
	}

	apiErr := api.SendDocument(&model.TelegramOutDocument{
		ChatId:    recipientId,
		Caption:   "*" + m.Origin + "* sent a document",
		ParseMode: "Markdown",
		Document: &model.TelegramInputFile{
			Name: m.Filename,
			Data: decodedFile,
		},
	})

	if apiErr != nil {
		log.Printf("error: %v\n", apiErr)
	}

	return apiErr
}
