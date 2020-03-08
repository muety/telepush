package internal

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"github.com/n1try/telegram-middleman-bot/model"
	"io"
	"mime/multipart"
	"net/http"
)

func validateFile(m *model.InMessage) error {
	if len(m.File) == 0 || len(m.Filename) == 0 {
		return errors.New("file or file name parameter missing")
	}
	return nil
}

func logFile(m *model.InMessage) string {
	return "A document named " + m.Filename + " was sent"
}

func resolveFile(recipientId string, m *model.InMessage) error {
	decodedFile, _ := b64.StdEncoding.DecodeString(m.File)
	return sendFile(recipientId, decodedFile, m.Filename, m.Origin)
}

func sendFile(recipientId string, file []byte, filename string, origin string) error {
	body := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(body)
	params := map[string]string{
		"caption":    "*" + origin + "* sent a document",
		"parse_mode": "Markdown",
	}
	addFileAndParamsToMultipartWriter(multipartWriter, file, filename, params)
	req, err := http.NewRequest("POST", botConfig.GetApiUrl()+"/sendDocument?chat_id="+recipientId, body)
	req.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}
	return nil
}

func addFileAndParamsToMultipartWriter(multipartWriter *multipart.Writer, file []byte, filename string, params map[string]string) {
	reader := bytes.NewReader(file)
	part, _ := multipartWriter.CreateFormFile("document", filename)
	io.Copy(part, reader)
	for key, val := range params {
		multipartWriter.WriteField(key, val)
	}
	multipartWriter.Close()
}
