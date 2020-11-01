package config

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

const (
	BaseURL          = "https://api.telegram.org/bot"
	StoreFile        = "store.gob"
	StatsFile        = "stats.json"
	PollTimeoutSec   = 60
	FlushTimeoutMin  = 60
	LimitsTimeoutMin = 60
)

const (
	KeyUpdateID = "latestUpdateId"
	KeyRequests = "totalRequests"
	KeyMessage  = "message"
	KeyParams   = "message_params"
	KeyMessages = "messages"
)

var cfg *BotConfig

type BotConfig struct {
	Token     string
	Mode      string
	UseHTTPS  bool
	CertPath  string
	KeyPath   string
	ProxyURI  *url.URL
	Port      int
	RateLimit int
	Address   string
	Address6  string
	Disable6  bool
	Version   string
}

func readVersion() string {
	file, err := os.Open("version.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func Get() *BotConfig {
	if cfg == nil {
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

		proxyUri, err := url.Parse(*proxyPtr)
		if err != nil || proxyUri.String() == "" {
			log.Println("Failed to parse proxy URI.")
		}

		cfg = &BotConfig{
			Token:     *tokenPtr,
			Mode:      *modePtr,
			UseHTTPS:  *useHttpsPtr,
			CertPath:  *certPathPtr,
			KeyPath:   *keyPathPtr,
			Port:      *portPtr,
			ProxyURI:  proxyUri,
			RateLimit: *rateLimitPtr,
			Address:   *addrPtr,
			Address6:  *addr6Ptr,
			Disable6:  *disable6Ptr,
			Version:   readVersion(),
		}
	}

	return cfg
}

func (c *BotConfig) GetApiUrl() string {
	return BaseURL + c.Token
}
