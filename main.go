package main

import (
	"fmt"
	"github.com/muety/telepush/handlers"
	"github.com/muety/telepush/inlets/alertmanager_webhook"
	"github.com/muety/telepush/inlets/bitbucket_webhook"
	"github.com/muety/telepush/inlets/webmentionio_webhook"
	"github.com/muety/telepush/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"github.com/muety/telepush/api"
	"github.com/muety/telepush/config"
	"github.com/muety/telepush/inlets/default"
	"github.com/muety/telepush/middleware"
)

var botConfig *config.BotConfig

func init() {
	botConfig = config.Get()
}

func registerRoutes() {
	indexHandler := handlers.NewIndexHandler()
	messageHandler := handlers.NewMessageHandler(
		services.NewUserService(config.GetStore()),
	)
	baseChain := alice.New(
		middleware.WithEventLogging(),
		middleware.WithMethodCheck(),
		middleware.WithRateLimit(),
	)

	if botConfig.Metrics {
		http.Handle("/metrics", promhttp.Handler())
	}
	http.Handle("/api/messages", baseChain.Append(_default.New().Handler).Then(messageHandler))
	http.Handle("/api/inlets/default", baseChain.Append(_default.New().Handler).Then(messageHandler))
	http.Handle("/api/inlets/alertmanager_webhook", baseChain.Append(alertmanager_webhook.New().Handler).Then(messageHandler))
	http.Handle("/api/inlets/bitbucket_webhook", baseChain.Append(bitbucket_webhook.New().Handler).Then(messageHandler))
	http.Handle("/api/inlets/webmentionio_webhook", baseChain.Append(webmentionio_webhook.New().Handler).Then(messageHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static/"))))
	http.Handle("/", indexHandler)
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
		fmt.Printf("Listening for HTTPS on %s.\n", s.Addr)
		go func() {
			if err := s.ListenAndServeTLS(botConfig.CertPath, botConfig.KeyPath); err != nil {
				log.Fatalln(err)
			}
		}()

		if s6 != nil {
			fmt.Printf("Listening for HTTPS on %s.\n", s6.Addr)
			go func() {
				if err := s6.ListenAndServeTLS(botConfig.CertPath, botConfig.KeyPath); err != nil {
					log.Fatalln(err)
				}
			}()
		}
	} else {
		fmt.Printf("Listening for HTTP on %s.\n", s.Addr)
		go func() {
			if err := s.ListenAndServe(); err != nil {
				log.Fatalln(err)
			}
		}()

		if s6 != nil {
			fmt.Printf("Listening for HTTP on %s.\n", s6.Addr)
			go func() {
				if err := s6.ListenAndServe(); err != nil {
					log.Fatalln(err)
				}
			}()
		}
	}
}

func exitGracefully() {
	config.GetHub().Close()
	config.GetStore().Flush()
}

func main() {
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
