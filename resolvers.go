package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	errors "errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

const TEXT_TYPE = "TEXT"
const FILE_TYPE = "FILE"

var typesResolvers map[string]TypeResolver

func InitResolvers() {
	typesResolvers = map[string]TypeResolver{
		TEXT_TYPE: TypeResolver{
			Resolve: resolveText,
			IsValid: validateText,
			Value:   logText,
		},
		FILE_TYPE: TypeResolver{
			Resolve: resolveFile,
			IsValid: validateFile,
			Value:   logFile,
		},
	}
}

// Send Message validaiton and resolving
func validateText(m InMessage) error {
	if len(m.Text) == 0 {
		return errors.New("You need to pass a text parameter")
	}
	return nil
}

func logText(m InMessage) string {
	return m.Text
}

func resolveText(recipientId string, m InMessage) error {
	return sendMessage(recipientId, "*"+m.Origin+"* wrote:\n\n"+m.Text)
}

func sendMessage(recipientId, text string) error {
	m, err := json.Marshal(&TelegramOutMessage{ChatId: recipientId, Text: text, ParseMode: "Markdown"})
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(m))
	request, _ := http.NewRequest(http.MethodPost, getApiUrl()+"/sendMessage", reader)
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Send Document validation and resolving
func validateFile(m InMessage) error {
	if len(m.File) == 0 || len(m.Filename) == 0 {
		return errors.New("You need to pass a file and filename parameter")
	}
	return nil
}

func logFile(m InMessage) string {
	return "A document named " + m.Filename + " was sent"
}

func resolveFile(recipientId string, m InMessage) error {
	decodedFile := decodeB64StringToByteArray(m.File)
	return sendFile(recipientId, decodedFile, m.Filename, m.Origin)
}

func decodeB64StringToByteArray(encodedString string) []byte {
	decodedString, _ := b64.StdEncoding.DecodeString(encodedString)
	return decodedString
}

func sendFile(recipientId string, file []byte, filename string, origin string) error {
	body := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(body)
	params := map[string]string{
		"caption":    "*" + origin + "* sent a document",
		"parse_mode": "Markdown",
	}
	addFileAndParamsToMultipartWriter(multipartWriter, file, filename, params)
	req, err := http.NewRequest("POST", getApiUrl()+"/sendDocument?chat_id="+recipientId, body)
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
