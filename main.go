package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

const BOT_API_TOKEN = "439545547:AAFBkRktGKTYnXKY-7Zr5TMIwF9RO1fl43M"
const BASE_URL = "https://api.telegram.org/bot"

//const BOT_API_TOKEN = ""
//const BASE_URL = "https://requestb.in/19mgpmo1"
const API_URL = BASE_URL + BOT_API_TOKEN
const STORE_FILE = "store.gob"
const POLL_TIMEOUT_SEC = 180
const STORE_KEY_UPDATE_ID = "latestUpdateId"
const STORE_KEY_REQUESTS = "totalRequests"

func sendMessage(recipientId, text string) error {
	m, err := json.Marshal(&TelegramOutMessage{ChatId: recipientId, Text: text, ParseMode: "Markdown"})
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(m))
	_, err = http.Post(API_URL+"/sendMessage", "application/json", reader)
	if err != nil {
		return err
	}
	return nil
}

func invalidateUserToken(userChatId int) {
	for k, v := range StoreGetMap() {
		entry, ok := v.(StoreObject)
		if ok && entry.ChatId == userChatId {
			StoreDelete(k)
		}
	}
}

func resolveToken(token string) string {
	value := StoreGet(token)
	if value != nil {
		return strconv.Itoa((value.(StoreObject)).ChatId)
	}
	return ""
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(415)
		return
	}
	StorePut(STORE_KEY_REQUESTS, StoreGet(STORE_KEY_REQUESTS).(int)+1)
	dec := json.NewDecoder(r.Body)
	var m InMessage
	err := dec.Decode(&m)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	if len(m.RecipientToken) == 0 || len(m.Text) == 0 || len(m.Origin) == 0 {
		w.Write([]byte("You need to pass recipient_token, origin and text parameters."))
		w.WriteHeader(400)
		return
	}

	recipientId := resolveToken(m.RecipientToken)
	if len(recipientId) == 0 {
		w.Write([]byte("The token you passed doesn't seem to relate to a valid user."))
		w.WriteHeader(404)
		return
	}

	err = sendMessage(recipientId, "*"+m.Origin+"* wrote:\n\n"+m.Text)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func getUpdate() (*[]TelegramUpdate, error) {
	offset := 0
	if StoreGet(STORE_KEY_UPDATE_ID) != nil {
		offset = int(StoreGet(STORE_KEY_UPDATE_ID).(float64)) + 1
	}
	url := API_URL + string("/getUpdates?timeout="+strconv.Itoa(POLL_TIMEOUT_SEC)+"&offset="+strconv.Itoa(offset))
	log.Println("Polling for updates.")
	request, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{Timeout: (POLL_TIMEOUT_SEC + 10) * time.Second}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var update TelegramUpdateResponse
	err = json.Unmarshal(data, &update)
	if err != nil {
		return nil, err
	}

	if len(update.Result) > 0 {
		var latestUpdateId interface{} = float64(update.Result[len(update.Result)-1].UpdateId)
		StorePut(STORE_KEY_UPDATE_ID, latestUpdateId)
	}
	return &update.Result, nil
}

func startPolling() {
	for {
		updates, err := getUpdate()
		if err == nil {
			for _, update := range *updates {
				var text string
				chatId := update.Message.Chat.Id
				if strings.HasPrefix(update.Message.Text, "/start") {
					id := uuid.NewV4().String()
					invalidateUserToken(chatId)
					StorePut(id, StoreObject{User: update.Message.From, ChatId: chatId})
					text = "Here is your token you can use to send messages to your Telegram account:\n\n_" + id + "_"
					log.Println("Sending new token to", strconv.Itoa(chatId))
				} else {
					text = "Please use the _/start_ command to fetch a new token.\n\nFurther information at https://github.com/n1try/telegram-middleman-bot."
				}
				err = sendMessage(strconv.Itoa(chatId), text)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func toJson(filePath string, data interface{}) {
	log.Println("Saving json.")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitStore()
	ReadStoreFromBinary(STORE_FILE)
	if StoreGet(STORE_KEY_REQUESTS) == nil {
		StorePut(STORE_KEY_REQUESTS, 0)
	}
	go startPolling()
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			FlushStoreToBinary(STORE_FILE)
			stats := Stats{TotalRequests: StoreGet(STORE_KEY_REQUESTS).(int), Timestamp: int(time.Now().Unix())}
			toJson("stats.json", stats)
		}
	}()

	// Exit handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	go func() {
		for _ = range c {
			FlushStoreToBinary(STORE_FILE)
			os.Exit(0)
		}
	}()

	http.HandleFunc("/api/messages", messageHandler)
	http.ListenAndServe(":8080", nil)
}
