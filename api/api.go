package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/muety/webhook2telegram/config"
	"github.com/muety/webhook2telegram/model"
	"github.com/muety/webhook2telegram/store"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	botConfig *config.BotConfig
	client    *http.Client
)

func init() {
	botConfig = config.Get()
	client = &http.Client{Timeout: (config.PollTimeoutSec + 10) * time.Second}
	if botConfig.ProxyURI != nil && botConfig.ProxyURI.String() != "" {
		client.Transport = &http.Transport{Proxy: http.ProxyURL(botConfig.ProxyURI)}
	}
}

func GetUpdate() (*[]model.TelegramUpdate, error) {
	offset := 0
	if store.Get(config.KeyUpdateID) != nil {
		offset = int(store.Get(config.KeyUpdateID).(float64)) + 1
	}
	apiUrl := botConfig.GetApiUrl() + string("/getUpdates?timeout="+strconv.Itoa(config.PollTimeoutSec)+"&offset="+strconv.Itoa(offset))
	log.Println("Polling for updates.")
	request, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(string(data))
	}

	var update model.TelegramUpdateResponse
	if err := json.Unmarshal(data, &update); err != nil {
		return nil, err
	}

	if len(update.Result) > 0 {
		var latestUpdateId interface{} = float64(update.Result[len(update.Result)-1].UpdateId)
		store.Put(config.KeyUpdateID, latestUpdateId)
	}

	return &update.Result, nil
}

func Poll() {
	go func() {
		for {
			if updates, err := GetUpdate(); err == nil {
				for _, update := range *updates {
					if err := processUpdate(update); err != nil {
						log.Printf("error processing updates: %s\n", err.Error())
					}
				}
			} else {
				log.Printf("error getting updates: %s\n", err)
				time.Sleep(config.PollTimeoutSec * time.Second)
			}
		}
	}()
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var u model.TelegramUpdate
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := processUpdate(u); err != nil {
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SendMessage(message *model.TelegramOutMessage) *model.ApiError {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(message); err != nil {
		return &model.ApiError{
			StatusCode: http.StatusBadRequest,
			Text:       err.Error(),
		}
	}
	request, _ := http.NewRequest(http.MethodPost, botConfig.GetApiUrl()+"/sendMessage", &buf)
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return &model.ApiError{
			StatusCode: http.StatusInternalServerError,
			Text:       err.Error(),
		}
	}
	defer resp.Body.Close()

	return handleApiResponse(resp)
}

func SendDocument(document *model.TelegramOutDocument) *model.ApiError {
	buf, contentType, err := document.EncodeMultipart()
	if err != nil {
		return &model.ApiError{
			StatusCode: http.StatusBadRequest,
			Text:       err.Error(),
		}
	}

	request, _ := http.NewRequest(http.MethodPost, botConfig.GetApiUrl()+"/sendDocument", buf)
	request.Header.Set("Content-Type", contentType)

	resp, err := client.Do(request)
	if err != nil {
		return &model.ApiError{
			StatusCode: http.StatusInternalServerError,
			Text:       err.Error(),
		}
	}
	defer resp.Body.Close()

	return handleApiResponse(resp)
}

func processUpdate(update model.TelegramUpdate) *model.ApiError {
	text := config.MessageDefaultResponse
	chatId := update.Message.Chat.Id

	if strings.HasPrefix(update.Message.Text, config.CmdStart) {
		id := uuid.NewV4()
		store.InvalidateToken(chatId)
		store.Put(id.String(), model.StoreObject{User: update.Message.From, ChatId: chatId})
		text = fmt.Sprintf(config.MessageTokenResponse, id.String())
		log.Printf("Sending new token %s to %s", id.String(), strconv.Itoa(chatId))
	} else if strings.HasPrefix(update.Message.Text, config.CmdHelp) {
		text = fmt.Sprintf(config.MessageHelpResponse, botConfig.Version)
	}

	return SendMessage(&model.TelegramOutMessage{
		ChatId:             strconv.Itoa(chatId),
		Text:               text,
		ParseMode:          "Markdown",
		DisableLinkPreview: true,
	})
}

func handleApiResponse(response *http.Response) *model.ApiError {
	resData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &model.ApiError{
			StatusCode: http.StatusInternalServerError,
			Text:       err.Error(),
		}
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(resData, &jsonResponse); err != nil {
		return &model.ApiError{
			StatusCode: http.StatusInternalServerError,
			Text:       err.Error(),
		}
	} else if ok := jsonResponse["ok"]; !(ok.(bool)) {
		desc := jsonResponse["description"].(string)
		status := jsonResponse["error_code"].(float64)
		return &model.ApiError{
			StatusCode: int(status),
			Text:       desc,
		}
	}

	return nil
}
