package main

import (
	"fmt"
	"github.com/n1try/telegram-middleman-bot/inlets/alertmanager_webhook"
	"github.com/n1try/telegram-middleman-bot/inlets/bitbucket_webhook"
	"github.com/n1try/telegram-middleman-bot/inlets/webmentionio_webhook"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"github.com/n1try/telegram-middleman-bot/api"
	"github.com/n1try/telegram-middleman-bot/config"
	"github.com/n1try/telegram-middleman-bot/inlets/default"
	"github.com/n1try/telegram-middleman-bot/middleware"
	"github.com/n1try/telegram-middleman-bot/model"
	"github.com/n1try/telegram-middleman-bot/resolvers"
	"github.com/n1try/telegram-middleman-bot/store"
	"github.com/n1try/telegram-middleman-bot/util"
)

var (
	botConfig  *config.BotConfig
	limiterMap map[string]int
)

type MessageHandler struct{}

func (h MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m *model.DefaultMessage
	var p *model.MessageParams

	if message := r.Context().Value(config.KeyMessage); message != nil {
		m = message.(*model.DefaultMessage)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse message"))
		return
	}

	if params := r.Context().Value(config.KeyParams); params != nil {
		p = params.(*model.MessageParams)
	}

	token := r.Header.Get("token")
	if token == "" {
		token = m.RecipientToken
	}

	if len(token) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing recipient_token parameter"))
		return
	}

	// TODO: Refactoring: get rid of this resolver concept
	resolver := resolvers.GetResolver(m.Type)

	if err := resolver.IsValid(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	recipientId := store.ResolveToken(token)

	if len(recipientId) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("passed token does not relate to a valid user"))
		return
	}

	_, hasKey := limiterMap[recipientId]
	if !hasKey {
		limiterMap[recipientId] = 0
	}
	if limiterMap[recipientId] >= botConfig.RateLimit {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(fmt.Sprintf("request rate of %d per hour exceeded", botConfig.RateLimit)))
		return
	}
	limiterMap[recipientId] += 1

	if err := resolver.Resolve(recipientId, m, p); err != nil {
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	store.Put(config.KeyRequests, store.Get(config.KeyRequests).(int)+1)

	w.WriteHeader(http.StatusOK)
}

func init() {
	botConfig = config.Get()
}

func flush() {
	for {
		time.Sleep(config.FlushTimeoutMin * time.Minute)
		store.Flush(config.StoreFile)

		stats := model.Stats{TotalRequests: store.Get(config.KeyRequests).(int), Timestamp: int(time.Now().Unix())}
		util.DumpJson(config.StatsFile, stats)
	}
}

func updateLimits() {
	for {
		limiterMap = make(map[string]int)
		time.Sleep(config.LimitsTimeoutMin * time.Minute)
	}
}

func registerRoutes() {
	messageHandler := &MessageHandler{}
	baseChain := alice.New(
		middleware.NewCheckMethod(botConfig),
		middleware.NewRateLimit(botConfig),
	)

	http.Handle("/api/messages", baseChain.Append(_default.New().Handler).Then(messageHandler))
	http.Handle("/api/inlets/default", baseChain.Append(_default.New().Handler).Then(messageHandler))
	http.Handle("/api/inlets/alertmanager_webhook", baseChain.Append(alertmanager_webhook.New().Handler).Then(messageHandler))
	http.Handle("/api/inlets/bitbucket_webhook", baseChain.Append(bitbucket_webhook.New().Handler).Then(messageHandler))
	http.Handle("/api/inlets/webmentionio_webhook", baseChain.Append(webmentionio_webhook.New().Handler).Then(messageHandler))
}

func connectApi() {
	if botConfig.Mode == "webhook" {
		fmt.Println("Using webhook mode.")
		http.HandleFunc("/api/updates", api.Webhook)
	} else {
		fmt.Println("Using long-polling mode.")
		api.Poll()
	}
}

func listen() {
	// Check if address is valid
	ip := net.ParseIP(botConfig.Address)
	if ip == nil {
		log.Println("Address '" + botConfig.Address + "' is not valid. Exiting...")
		os.Exit(1)
	}

	// IPv4
	bindString := botConfig.Address + ":" + strconv.Itoa(botConfig.Port)
	s := &http.Server{
		Addr:         bindString,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// IPv6
	var s6 *http.Server
	if !botConfig.Disable6 {
		ip := net.ParseIP(botConfig.Address6)
		if ip == nil {
			log.Println("Address '" + botConfig.Address6 + "' is not valid. Exiting...")
			os.Exit(1)
		}

		bindString := "[" + botConfig.Address6 + "]:" + strconv.Itoa(botConfig.Port)
		s6 = &http.Server{
			Addr:         bindString,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
	}

	if botConfig.UseHTTPS {
		fmt.Printf("Listening for HTTPS on port %d.\n", botConfig.Port)
		if !botConfig.Disable6 {
			go s6.ListenAndServeTLS(botConfig.CertPath, botConfig.KeyPath)
		}
		go s.ListenAndServeTLS(botConfig.CertPath, botConfig.KeyPath)
	} else {
		fmt.Printf("Listening for HTTP on port %d.\n", botConfig.Port)
		if !botConfig.Disable6 {
			go s6.ListenAndServe()
		}
		go s.ListenAndServe()
	}
}

func exitGracefully() {
	store.Flush(config.StoreFile)
}

func main() {
	store.Read(config.StoreFile)
	store.Automigrate()

	// Stats
	if store.Get(config.KeyRequests) == nil {
		store.Put(config.KeyRequests, 0)
	}

	go flush()
	go updateLimits()

	registerRoutes()
	connectApi()
	listen()

	// Exit handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Kill)

	<-c
	exitGracefully()
}
