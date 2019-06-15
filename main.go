package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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

const BASE_URL = "https://api.telegram.org/bot"
const STORE_FILE = "store.gob"
const POLL_TIMEOUT_SEC = 180
const STORE_KEY_UPDATE_ID = "latestUpdateId"
const STORE_KEY_REQUESTS = "totalRequests"
const STORE_KEY_MESSAGES = "messages"

var token string
var limiterMap map[string]int
var maxReqsPerHous int

func getApiUrl() string {
	return BASE_URL + token
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

	dec := json.NewDecoder(r.Body)
	var m InMessage
	err := dec.Decode(&m)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	if len(m.RecipientToken) == 0 || len(m.Origin) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to pass recipient_token and origin."))
		return
	}

	messageType := TEXT_TYPE
	if len(m.Type) > 0 {
		messageType = m.Type
	}

	invalid := typesResolvers[messageType].IsValid(m)
	if invalid != nil {
		w.WriteHeader(400)
		w.Write([]byte(invalid.Error()))
	}

	recipientId := resolveToken(m.RecipientToken)
	if len(recipientId) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("The token you passed doesn't seem to relate to a valid user."))
		return
	}

	_, hasKey := limiterMap[recipientId]
	if !hasKey {
		limiterMap[recipientId] = 0
	}
	if limiterMap[recipientId] >= maxReqsPerHous {
		w.WriteHeader(429)
		w.Write([]byte(fmt.Sprintf("Request rate of %d per hour exceeded.", maxReqsPerHous)))
		return
	}
	limiterMap[recipientId] += 1

	err = typesResolvers[messageType].Resolve(recipientId, m)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	valueForStore  := typesResolvers[messageType].ValueForStore(m)
	storedMessages := StoreGet(STORE_KEY_MESSAGES).(StoreMessageObject)
	storedMessages = append(storedMessages, valueForStore)
	StorePut(STORE_KEY_MESSAGES, storedMessages)
	StorePut(STORE_KEY_REQUESTS, StoreGet(STORE_KEY_REQUESTS).(int)+1)

	w.WriteHeader(200)
}

func webhookUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(415)
		return
	}
	dec := json.NewDecoder(r.Body)
	var u TelegramUpdate
	err := dec.Decode(&u)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	processUpdate(u)
	w.WriteHeader(200)
}

func getUpdate() (*[]TelegramUpdate, error) {
	offset := 0
	if StoreGet(STORE_KEY_UPDATE_ID) != nil {
		offset = int(StoreGet(STORE_KEY_UPDATE_ID).(float64)) + 1
	}
	url := getApiUrl() + string("/getUpdates?timeout="+strconv.Itoa(POLL_TIMEOUT_SEC)+"&offset="+strconv.Itoa(offset))
	log.Println("Polling for updates.")
	request, _ := http.NewRequest("GET", url, nil)
	request.Close = true
	client := &http.Client{Timeout: (POLL_TIMEOUT_SEC + 10) * time.Second}

	response, err := client.Do(request)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(string(data))
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

func processUpdate(update TelegramUpdate) {
	var text string
	chatId := update.Message.Chat.Id
	if strings.HasPrefix(update.Message.Text, "/start") {
		u, _ := uuid.NewV4()
		id := u.String()
		invalidateUserToken(chatId)
		StorePut(id, StoreObject{User: update.Message.From, ChatId: chatId})
		text = "Here is your token you can use to send messages to your Telegram account:\n\n_" + id + "_"
		log.Println("Sending new token to", strconv.Itoa(chatId))
	} else {
		text = "Please use the _/start_ command to fetch a new token.\n\nFurther information at https://github.com/n1try/telegram-middleman-bot."
	}
	err := sendMessage(strconv.Itoa(chatId), text)
	if err != nil {
		log.Println(err)
	}
}

func startPolling() {
	for {
		updates, err := getUpdate()
		if err == nil {
			for _, update := range *updates {
				processUpdate(update)
			}
		} else {
			log.Printf("ERROR getting updates: %s\n", err)
			time.Sleep(POLL_TIMEOUT_SEC * time.Second)
		}
	}
}

func getConfig() BotConfig {
	tokenPtr := flag.String("token", "", "Your Telegram Bot Token from Botfather")
	modePtr := flag.String("mode", "poll", "Update mode ('poll' for development, 'webhook' for production)")
	useHttpsPtr := flag.Bool("useHttps", false, "Whether or not to use TLS for webserver. Required for webhook mode if not using a reverse proxy")
	certPathPtr := flag.String("certPath", "", "Path of your SSL certificate when using webhook mode")
	keyPathPtr := flag.String("keyPath", "", "Path of your private SSL key when using webhook mode")
	portPtr := flag.Int("port", 8080, "Port for the webserver to listen on")
	rateLimitPtr := flag.Int("rateLimit", 10, "Max number of requests per recipient per hour")

	flag.Parse()

	return BotConfig{
		Token:     *tokenPtr,
		Mode:      *modePtr,
		UseHTTPS:  *useHttpsPtr,
		CertPath:  *certPathPtr,
		KeyPath:   *keyPathPtr,
		Port:      *portPtr,
		RateLimit: *rateLimitPtr}
}

func toJson(filePath string, data interface{}) {
	log.Println("Saving json.")
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&data)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	InitStore()
	InitResolvers()
	ReadStoreFromBinary(STORE_FILE)

	go func() {
		for {
			time.Sleep(60 * time.Minute)
			FlushStoreToBinary(STORE_FILE)
			stats := Stats{TotalRequests: StoreGet(STORE_KEY_REQUESTS).(int), Timestamp: int(time.Now().Unix())}
			toJson("stats.json", stats)

			limiterMap = make(map[string]int)
		}
	}()

	config := getConfig()
	token = config.Token
	limiterMap = make(map[string]int)
	maxReqsPerHous = config.RateLimit

	http.HandleFunc("/api/messages", messageHandler)

	if config.Mode == "webhook" {
		fmt.Println("Using webhook mode.")
		http.HandleFunc("/api/updates", webhookUpdateHandler)
	} else {
		fmt.Println("Using long-polling mode.")
		go startPolling()
	}

	if StoreGet(STORE_KEY_REQUESTS) == nil {
		StorePut(STORE_KEY_REQUESTS, 0)
	}

	if StoreGet(STORE_KEY_MESSAGES) == nil {
		StorePut(STORE_KEY_MESSAGES, StoreMessageObject{})
	}

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

	portString := ":" + strconv.Itoa(config.Port)
	s := &http.Server{
		Addr:         portString,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if config.UseHTTPS {
		fmt.Printf("Listening for HTTPS on port %d.\n", config.Port)
		s.ListenAndServeTLS(config.CertPath, config.KeyPath)
	} else {
		fmt.Printf("Listening for HTTP on port %d.\n", config.Port)
		s.ListenAndServe()
	}
}
