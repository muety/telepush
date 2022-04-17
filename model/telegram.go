package model

import (
	"bytes"
	"io"
	"mime/multipart"
)

// Only required fields are implemented
type TelegramUser struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Only required fields are implemented
type TelegramChat struct {
	Id   int64  `json:"id"`
	Type string `json:"type"`
}

// Only required fields are implemented
type TelegramOutMessage struct {
	ChatId             string `json:"chat_id"`
	Text               string `json:"text"`
	ParseMode          string `json:"parse_mode"`
	DisableLinkPreview bool   `json:"disable_web_page_preview"`
}

type TelegramOutDocument struct {
	ChatId    string
	Caption   string
	ParseMode string
	Document  *TelegramInputFile
}

type TelegramInputFile struct {
	Name string
	Data []byte
}

// Only required fields are implemented
type TelegramInMessage struct {
	MessageId int64        `json:"message_id"`
	From      TelegramUser `json:"from"`
	Date      int64        `json:"date"`
	Chat      TelegramChat `json:"chat"`
	Text      string       `json:"text"`
}

// Only required fields are implemented
type TelegramUpdate struct {
	UpdateId int64             `json:"update_id"`
	Message  TelegramInMessage `json:"message"`
}

type TelegramUpdateResponse struct {
	Ok     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

func (d *TelegramOutDocument) EncodeMultipart() (*bytes.Buffer, string, error) {
	buf := bytes.Buffer{}
	w := multipart.NewWriter(&buf)
	defer w.Close()

	filePart, err := w.CreateFormFile("document", d.Document.Name)
	if err != nil {
		return nil, "", err
	}

	if _, err := io.Copy(filePart, bytes.NewReader(d.Document.Data)); err != nil {
		return nil, "", err
	}

	w.WriteField("chat_id", d.ChatId)
	w.WriteField("caption", d.Caption)
	w.WriteField("parse_mode", d.ParseMode)

	return &buf, w.FormDataContentType(), nil
}
