package main

import (
	"fmt"
	"github.com/n1try/telegram-middleman-bot/handlers"
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
	"github.com/n1try/telegram-middleman-bot/store"
	"github.com/n1try/telegram-middleman-bot/util"
)

var (
	botConfig *config.BotConfig
)

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
		time.Sleep(config.LimitsTimeoutMin * time.Minute)
	}
}

func registerRoutes() {
	messageHandler := &handlers.MessageHandler{}
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
		log.Println("address '" + botConfig.Address + "' is not valid. Exiting...")
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
			log.Println("address '" + botConfig.Address6 + "' is not valid. Exiting...")
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
