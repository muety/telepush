package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	BaseURL          = "https://api.telegram.org/bot"
	StoreFile        = "store.gob"
	PollTimeoutSec   = 60
	StoreKeyUpdateID = "latestUpdateId"
	StoreKeyRequests = "totalRequests"
	StoreKeyMessages = "messages"
)

var (
	token          string
	limiterMap     map[string]int
	maxReqsPerHour int
	client         = &http.Client{Timeout: (PollTimeoutSec + 10) * time.Second}
)

func getApiUrl() string {
	return BaseURL + token
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

	token := r.Header.Get("token")
	if token == "" {
		token = m.RecipientToken
	}

	if len(token) == 0 || len(m.Origin) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to pass recipient_token and origin."))
		return
	}

	messageType := TextType
	if len(m.Type) > 0 {
		messageType = m.Type
	}

	invalid := typesResolvers[messageType].IsValid(m)
	if invalid != nil {
		w.WriteHeader(400)
		w.Write([]byte(invalid.Error()))
		return
	}

	recipientId := resolveToken(token)

	if len(recipientId) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("The token you passed doesn't seem to relate to a valid user."))
		return
	}

	_, hasKey := limiterMap[recipientId]
	if !hasKey {
		limiterMap[recipientId] = 0
	}
	if limiterMap[recipientId] >= maxReqsPerHour {
		w.WriteHeader(429)
		w.Write([]byte(fmt.Sprintf("Request rate of %d per hour exceeded.", maxReqsPerHour)))
		return
	}
	limiterMap[recipientId] += 1

	err = typesResolvers[messageType].Resolve(recipientId, m)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	valueForStore := typesResolvers[messageType].Value(m)
	storedMessages := StoreGet(StoreKeyMessages).(StoreMessageObject)
	storedMessages = append(storedMessages, valueForStore)
	StorePut(StoreKeyMessages, storedMessages)
	StorePut(StoreKeyRequests, StoreGet(StoreKeyRequests).(int)+1)

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
	if StoreGet(StoreKeyUpdateID) != nil {
		offset = int(StoreGet(StoreKeyUpdateID).(float64)) + 1
	}
	apiUrl := getApiUrl() + string("/getUpdates?timeout="+strconv.Itoa(PollTimeoutSec)+"&offset="+strconv.Itoa(offset))
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

	var update TelegramUpdateResponse
	err = json.Unmarshal(data, &update)
	if err != nil {
		return nil, err
	}

	if len(update.Result) > 0 {
		var latestUpdateId interface{} = float64(update.Result[len(update.Result)-1].UpdateId)
		StorePut(StoreKeyUpdateID, latestUpdateId)
	}

	return &update.Result, nil
}

func processUpdate(update TelegramUpdate) {
	var text string
	chatId := update.Message.Chat.Id
	if strings.HasPrefix(update.Message.Text, "/start") {
		id, _ := uuid.NewV4()

		invalidateUserToken(chatId)
		StorePut(id.String(), StoreObject{User: update.Message.From, ChatId: chatId})
		text = "Here is your token you can use to send messages to your Telegram account:\n\n_" + id.String() + "_"
		log.Printf("Sending new token %s to %s", id.String(), strconv.Itoa(chatId))
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
			time.Sleep(PollTimeoutSec * time.Second)
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
	proxyPtr := flag.String("proxy", "", "Proxy for poll mode, e.g. 'socks5://127.0.0.01:1080'")
	rateLimitPtr := flag.Int("rateLimit", 10, "Max number of requests per recipient per hour")
	addrPtr := flag.String("address", "127.0.0.1", "IPv4 address to bind the webserver to")
	addr6Ptr := flag.String("address6", "::1", "IPv6 address to bind the webserver to")
	disable6Ptr := flag.Bool("disableIPv6", false, "Set if your device doesn't support IPv6. address6 will be ignored if this is set.")

	flag.Parse()

	return BotConfig{
		Token:     *tokenPtr,
		Mode:      *modePtr,
		UseHTTPS:  *useHttpsPtr,
		CertPath:  *certPathPtr,
		KeyPath:   *keyPathPtr,
		Port:      *portPtr,
		ProxyURI:  *proxyPtr,
		RateLimit: *rateLimitPtr,
		Address:   *addrPtr,
		Address6:  *addr6Ptr,
		Disable6:  *disable6Ptr}
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
	ReadStoreFromBinary(StoreFile)

	go func() {
		for {
			time.Sleep(60 * time.Minute)
			FlushStoreToBinary(StoreFile)
			stats := Stats{TotalRequests: StoreGet(StoreKeyRequests).(int), Timestamp: int(time.Now().Unix())}
			toJson("stats.json", stats)

			limiterMap = make(map[string]int)
		}
	}()

	config := getConfig()

	if urlProxy, err := url.Parse(config.ProxyURI); err == nil && urlProxy.String() != "" {
		client.Transport = &http.Transport{Proxy: http.ProxyURL(urlProxy)}
	}

	token = config.Token
	limiterMap = make(map[string]int)
	maxReqsPerHour = config.RateLimit

	http.HandleFunc("/api/messages", messageHandler)

	if config.Mode == "webhook" {
		fmt.Println("Using webhook mode.")
		http.HandleFunc("/api/updates", webhookUpdateHandler)
	} else {
		fmt.Println("Using long-polling mode.")
		go startPolling()
	}

	if StoreGet(StoreKeyRequests) == nil {
		StorePut(StoreKeyRequests, 0)
	}

	if StoreGet(StoreKeyMessages) == nil {
		StorePut(StoreKeyMessages, StoreMessageObject{})
	}

	// Exit handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Kill)
	go func() {
		for _ = range c {
			FlushStoreToBinary(StoreFile)
			os.Exit(0)
		}
	}()

	// Check if address is valid
	ip := net.ParseIP(config.Address)
	if ip == nil {
		log.Println("Address '" + config.Address + "' is not valid. Exiting...")
		os.Exit(1)
	}

	// IPv4
	bindString := config.Address + ":" + strconv.Itoa(config.Port)
	s := &http.Server{
		Addr:         bindString,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// IPv6
	var s6 *http.Server
	if !config.Disable6 {
		ip := net.ParseIP(config.Address6)
		if ip == nil {
			log.Println("Address '" + config.Address6 + "' is not valid. Exiting...")
			os.Exit(1)
		}

		bindString := "[" + config.Address6 + "]:" + strconv.Itoa(config.Port)
		s6 = &http.Server{
			Addr:         bindString,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
	}

	if config.UseHTTPS {
		fmt.Printf("Listening for HTTPS on port %d.\n", config.Port)
		if !config.Disable6 {
			go s6.ListenAndServeTLS(config.CertPath, config.KeyPath)
		}
		s.ListenAndServeTLS(config.CertPath, config.KeyPath)
	} else {
		fmt.Printf("Listening for HTTP on port %d.\n", config.Port)
		if !config.Disable6 {
			go s6.ListenAndServe()
		}
		s.ListenAndServe()
	}
}
