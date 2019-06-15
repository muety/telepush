package main

import (
	"encoding/json"
	"net/http"
	"mime/multipart"
	"strings"
	b64 "encoding/base64"
	"bytes"
	"io"
	errors "errors"
)

const TEXT_TYPE = "TEXT"
const FILE_TYPE = "FILE"
var typesResolvers map[string] TypeResolver

func InitResolvers() {
	typesResolvers = map[string] TypeResolver {
		TEXT_TYPE: TypeResolver{
			Resolve:       sendMessageClosure,
			IsValid:       sendMessageValidation,
			ValueForStore: sendMessageValueForStore,
		},
		FILE_TYPE: TypeResolver{
			Resolve:       sendFileClosure,
			IsValid:       sendDocumentValidation,
			ValueForStore: sendDocumentValueForStore,
		},
	}
}

// Send Message validaiton and resolving
func sendMessageValidation(m InMessage) error {
	if len(m.Text) == 0 {
		return errors.New("You need to pass a text parameter")
	}
	return nil
}

func sendMessageValueForStore(m InMessage) string {
	return m.Text
}

func sendMessageClosure(recipientId string, m InMessage) error {
	return sendMessage(recipientId, "*"+m.Origin+"* wrote:\n\n"+m.Text)
}

func sendMessage(recipientId, text string) error {
	m, err := json.Marshal(&TelegramOutMessage{ChatId: recipientId, Text: text, ParseMode: "Markdown"})
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(m))
	resp, err := http.Post(getApiUrl()+"/sendMessage", "application/json", reader)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

// Send Document validation and resolving
func sendDocumentValidation(m InMessage) error {
	if len(m.File) == 0 || len(m.Filename) == 0 {
		return errors.New("You need to pass a file and filename parameter")
	}
	return nil
}

func sendDocumentValueForStore(m InMessage) string {
	return "A document named " + m.Filename + " was sent"
}

func sendFileClosure(recipientId string, m InMessage) error {
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
	params := map[string]string {
		"caption": "*" + origin + "* sent a document",
		"parse_mode": "Markdown",
	}
	addFileAndParamsToMultipartWriter(multipartWriter, file, filename, params)
	req, err := http.NewRequest("POST", getApiUrl()+"/sendDocument?chat_id=" + recipientId, body)
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
