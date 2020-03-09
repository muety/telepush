package api

import (
	"encoding/json"
	"errors"
	"github.com/n1try/telegram-middleman-bot/config"
	"github.com/n1try/telegram-middleman-bot/model"
	"github.com/n1try/telegram-middleman-bot/store"
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
	request, _ := http.NewRequest("GET", apiUrl, nil)
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

	if response.StatusCode != 200 {
		return nil, errors.New(string(data))
	}

	var update model.TelegramUpdateResponse
	err = json.Unmarshal(data, &update)
	if err != nil {
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
			updates, err := GetUpdate()
			if err == nil {
				for _, update := range *updates {
					processUpdate(update)
				}
			} else {
				log.Printf("ERROR getting updates: %s\n", err)
				time.Sleep(config.PollTimeoutSec * time.Second)
			}
		}
	}()
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(415)
		return
	}
	dec := json.NewDecoder(r.Body)
	var u model.TelegramUpdate
	err := dec.Decode(&u)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	processUpdate(u)
	w.WriteHeader(200)
}

func SendMessage(message *model.TelegramOutMessage) error {
	m, err := json.Marshal(message)
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(m))
	request, _ := http.NewRequest(http.MethodPost, botConfig.GetApiUrl()+"/sendMessage", reader)
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

func processUpdate(update model.TelegramUpdate) {
	var text string
	chatId := update.Message.Chat.Id
	if strings.HasPrefix(update.Message.Text, "/start") {
		id := uuid.NewV4()
		store.InvalidateToken(chatId)
		store.Put(id.String(), model.StoreObject{User: update.Message.From, ChatId: chatId})
		text = "Here is your token you can use to send messages to your Telegram account:\n\n_" + id.String() + "_"
		log.Printf("Sending new token %s to %s", id.String(), strconv.Itoa(chatId))
	} else {
		text = "Please use the _/start_ command to fetch a new token.\n\nFurther information at https://github.com/n1try/telegram-middleman-bot."
	}
	err := SendMessage(&model.TelegramOutMessage{
		ChatId:             strconv.Itoa(chatId),
		Text:               text,
		ParseMode:          "Markdown",
		DisableLinkPreview: true,
	})
	if err != nil {
		log.Println(err)
	}
}
