package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/muety/telepush/handlers"
	alertmanagerIn "github.com/muety/telepush/inlets/alertmanager"
	bitbucketIn "github.com/muety/telepush/inlets/bitbucket"
	grafanaIn "github.com/muety/telepush/inlets/grafana"
	webmentionioIn "github.com/muety/telepush/inlets/webmentionio"
	"github.com/muety/telepush/services"
	"github.com/muety/telepush/util"
	"github.com/muety/telepush/views"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/muety/telepush/api"
	"github.com/muety/telepush/config"
	defaultIn "github.com/muety/telepush/inlets/default"
	"github.com/muety/telepush/middleware"
)

var botConfig *config.BotConfig

func init() {
	botConfig = config.Get()
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
		log.Printf("Listening for HTTPS on %s.\n", s.Addr)
		go func() {
			if err := s.ListenAndServeTLS(botConfig.CertPath, botConfig.KeyPath); err != nil {
				log.Fatalln(err)
			}
		}()

		if s6 != nil {
			log.Printf("Listening for HTTPS on %s.\n", s6.Addr)
			go func() {
				if err := s6.ListenAndServeTLS(botConfig.CertPath, botConfig.KeyPath); err != nil {
					log.Fatalln(err)
				}
			}()
		}
	} else {
		log.Printf("Listening for HTTP on %s.\n", s.Addr)
		go func() {
			if err := s.ListenAndServe(); err != nil {
				log.Fatalln(err)
			}
		}()

		if s6 != nil {
			log.Printf("Listening for HTTP on %s.\n", s6.Addr)
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
	log.Printf("Environment: %s\n", botConfig.Env)
	log.Printf("Version: %s\n", botConfig.Version)

	// Initialize Router
	rootRouter := mux.NewRouter().StrictSlash(true)
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	apiRouter.Use(
		middleware.WithEventLogging(),
		middleware.WithRateLimit(),
	)

	// Initialize Handlers
	indexHandler := handlers.NewIndexHandler()
	messageHandler := handlers.NewMessageHandler(services.NewUserService(config.GetStore()))

	// Register Routes
	messageChain := alice.New(middleware.WithToken("recipient", config.KeyRecipient))
	apiRouter.Methods(http.MethodGet, http.MethodPost).Path("/messages/{recipient}").Handler(messageChain.Append(defaultIn.New().Handler).Then(messageHandler))
	apiRouter.Methods(http.MethodGet, http.MethodPost).Path("/inlets/default/{recipient}").Handler(messageChain.Append(defaultIn.New().Handler).Then(messageHandler))
	apiRouter.Methods(http.MethodGet, http.MethodPost).Path("/inlets/alertmanager/{recipient}").Handler(messageChain.Append(alertmanagerIn.New().Handler).Then(messageHandler))
	apiRouter.Methods(http.MethodGet, http.MethodPost).Path("/inlets/grafana/{recipient}").Handler(messageChain.Append(grafanaIn.New().Handler).Then(messageHandler))
	apiRouter.Methods(http.MethodGet, http.MethodPost).Path("/inlets/bitbucket/{recipient}").Handler(messageChain.Append(bitbucketIn.New().Handler).Then(messageHandler))
	apiRouter.Methods(http.MethodGet, http.MethodPost).Path("/inlets/webmentionio/{recipient}").Handler(messageChain.Append(webmentionioIn.New().Handler).Then(messageHandler))

	if botConfig.Mode == "webhook" {
		apiRouter.Methods(http.MethodPost).Path(botConfig.GetUpdatesPath()).HandlerFunc(api.Webhook)

		log.Println("Using webhook mode")
		log.Printf("Updates from Telegram are accepted under '/api%s'. Set the webhook accordingly (see https://core.telegram.org/bots/api#setwebhook)\n", botConfig.GetUpdatesPath())
		if botConfig.UrlSecret == "" {
			log.Println("Warning: It is recommended to set '-urlSecret' for enhanced security (can be any random string)")
		}
	} else {
		log.Println("Using long-polling mode")
		log.Println("Warning: For production use, webhook mode is recommended")
		api.Poll()
	}

	if botConfig.Metrics {
		rootRouter.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())
	}

	staticFs := util.NeuteredFileSystem{FS: views.GetStaticFilesFS()}
	rootRouter.Methods(http.MethodGet).PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.FS(staticFs))))
	rootRouter.Methods(http.MethodGet).PathPrefix("/").Handler(indexHandler)

	// Start server
	http.Handle("/", middleware.WithTrailingSlash()(rootRouter))
	listen()

	// Exit handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Kill)

	<-c
	exitGracefully()
}
