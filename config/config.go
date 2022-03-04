package config

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	BaseURL        = "https://api.telegram.org/bot"
	StoreFile      = "store.gob"
	PollTimeoutSec = 60
	UserIdRegex    = "(?m)^\\d+$"
)

const (
	KeyUpdateID  = "latestUpdateId"
	KeyMessage   = "message"
	KeyParams    = "message_params"
	KeyRecipient = "recipient"
)

const (
	CmdStart = "/start"
	CmdHelp  = "/help"
)

const (
	MessageDefaultResponse = "Please use the _/start_ command to fetch a new token.\n\nFurther information at https://github.com/muety/telepush."
	MessageTokenResponse   = "Here is your token, which you can include to request to this bot:\n\n`%s`"
	MessageHelpResponse    = "For detailed instructions on how to use this bot, please refer to the [official documentation](https://github.com/muety/telepush).\n\nVersion: `%s`"
)

var cfg *BotConfig

type BotConfig struct {
	Env          string
	Token        string
	Mode         string
	BaseUrl      string
	UseHTTPS     bool
	CertPath     string
	KeyPath      string
	ProxyURI     *url.URL
	Port         int
	ReqRateLimit int
	CmdRateLimit int
	Address      string
	Address6     string
	Disable6     bool
	Metrics      bool
	DataDir      string
	Blacklist    []string
	Version      string
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

func readBlacklist(path string) []string {
	if path == "" {
		return []string{}
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(UserIdRegex)
	lines := strings.Split(string(bytes), "\n")
	blacklist := make([]string, 0, len(lines))

	for _, l := range lines {
		if !re.MatchString(l) {
			continue
		}
		blacklist = append(blacklist, l)
	}

	return blacklist
}

func Get() *BotConfig {
	if cfg == nil {
		envPtr := flag.String("env", "production", "Environment to run in (dev or production)")
		tokenPtr := flag.String("token", "", "Your Telegram Bot Token from Botfather")
		modePtr := flag.String("mode", "poll", "Update mode ('poll' for development, 'webhook' for production)")
		baseUrlPtr := flag.String("baseUrl", "", "A relative URL different from '/', required to run the bot on a subpath. E.g. to run bot under 'https://exmaple.org/wh2tg' set baseUrl to '/wh2tg'")
		useHttpsPtr := flag.Bool("useHttps", false, "Whether or not to use TLS for webserver. Required for webhook mode if not using a reverse proxy")
		certPathPtr := flag.String("certPath", "", "Path of your SSL certificate when using webhook mode")
		keyPathPtr := flag.String("keyPath", "", "Path of your private SSL key when using webhook mode")
		portPtr := flag.Int("port", 8080, "Port for the webserver to listen on")
		proxyPtr := flag.String("proxy", "", "Proxy for poll mode, e.g. 'socks5://127.0.0.01:1080'")
		reqRateLimitPtr := flag.Int("rateLimit", 100, "Max number of requests per recipient per hour")
		cmdRateLimitPtr := flag.Int("cmdRateLimit", 10, "Max number of chat commands to execute per hour")
		addrPtr := flag.String("address", "127.0.0.1", "IPv4 address to bind the webserver to")
		addr6Ptr := flag.String("address6", "::1", "IPv6 address to bind the webserver to")
		disable6Ptr := flag.Bool("disableIPv6", false, "Set if your device doesn't support IPv6. address6 will be ignored if this is set.")
		metricsPtr := flag.Bool("metrics", false, "Whether or not to expose Prometheus metrics under '/metrics'")
		dataDirPtr := flag.String("dataDir", ".", "File system location where to store persistent data")
		blacklistPtr := flag.String("blacklist", "", "Path to a user id blacklist file (e.g. 'blacklist.txt')")

		flag.Parse()

		if *tokenPtr == "" {
			log.Fatalln("token missing")
		}

		proxyUri, err := url.Parse(*proxyPtr)
		if err != nil {
			log.Println("failed to parse proxy uri")
		}

		cfg = &BotConfig{
			Env:          *envPtr,
			Token:        *tokenPtr,
			Mode:         *modePtr,
			BaseUrl:      *baseUrlPtr + "/",
			UseHTTPS:     *useHttpsPtr,
			CertPath:     *certPathPtr,
			KeyPath:      *keyPathPtr,
			Port:         *portPtr,
			ProxyURI:     proxyUri,
			ReqRateLimit: *reqRateLimitPtr,
			CmdRateLimit: *cmdRateLimitPtr,
			Address:      *addrPtr,
			Address6:     *addr6Ptr,
			Disable6:     *disable6Ptr,
			Metrics:      *metricsPtr,
			DataDir:      *dataDirPtr,
			Blacklist:    readBlacklist(*blacklistPtr),
			Version:      readVersion(),
		}
	}

	return cfg
}

func (c *BotConfig) GetApiUrl() string {
	return BaseURL + c.Token
}

func (c *BotConfig) GetStorePath() string {
	return path.Join(c.DataDir, StoreFile)
}
