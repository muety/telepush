# webhook2telegram
(formerly _telegram-middleman-bot_)

[![](http://img.shields.io/liberapay/receives/muety.svg?logo=liberapay)](https://liberapay.com/muety/)
![](https://badges.fw-web.space/github/license/muety/webhook2telegram)
![Coding Activity](https://badges.fw-web.space/endpoint?url=https://wakapi.dev/api/compat/shields/v1/n1try/interval:any/project:webhook2telegram&color=blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/muety/webhook2telegram)](https://goreportcard.com/report/github.com/muety/webhook2telegram)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=muety_telegram-middleman-bot&metric=security_rating)](https://sonarcloud.io/dashboard?id=muety_telegram-middleman-bot)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=muety_telegram-middleman-bot&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=muety_telegram-middleman-bot)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=muety_telegram-middleman-bot&metric=sqale_index)](https://sonarcloud.io/dashboard?id=muety_telegram-middleman-bot)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=muety_telegram-middleman-bot&metric=ncloc)](https://sonarcloud.io/dashboard?id=muety_telegram-middleman-bot)

[![Buy me a coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://buymeacoff.ee/n1try)

---

![](views/static/logo.png)

A [Telegram Bot](https://telegram.me/MiddleManBot) to translate simple JSON HTTP requests into Telegram push messages that you will get on your Smartphone, PC or whatever Telegram client you have. Just like [Gotify](https://gotify.net/), but without an extra app.

## Changelog
### 2020-12-05
* Support for Prometheus metrics ([#18](https://github.com/muety/webhook2telegram/issues/18))
* Official Docker image ([`n1try/webhook2telegram`](https://hub.docker.com/repository/docker/n1try/webhook2telegram)) ([#20](https://github.com/muety/webhook2telegram/issues/20))

### 2020-11-01
* Project was renamed from _telegram-middleman-bot_ to _webhook2telegram_

### 2020-04-05
* Integration with [Webmention.io](https://webmention.io)

### 2020-03-26
* ⚠️ `alertmanager` inlet was renamed  to `alertmanager_webhook` and its endpoint has changed accordingly
* Docker support, thanks to [luza](https://github.com/luza)
* Integration with Bitbucket, thanks to [luza](https://github.com/luza)

### 2020-03-10
* Major code refactorings
* Support for Inlets (see below)
* Integration with Prometheus Alertmanager

### 2019-11-06
Thanks to contributions by [peet1993](https://github.com/peet1993).
* Introduced explicit IPv6 support 
* Introduced ability to specify network address to bind to

## Why might this be useful?
This is especially useful for __developers or sysadmins__. Imagine you want some kind of reporting from your application or server, like a daily report including some statistics. You don't want to actively look it up on a website but you want to receive it in a __passive fashion__. Just like getting an e-mail. But come on, let's be honest. __E-Mails are so 2010__. And they require your little server-side script to include some SMTP library and connect to a mail server. That's __too heavyweight__ just to __get some short information__. Personally, I have a Python script running on my server which gathers some statistics from log files and databases and regularly sends me a Telegram message.

If you develop those thoughts further, this could potentially __replace any kind of e-mail notifications__ - be it the message that someone has answered to your __forum post__, your favorite game is __now on sale at Steam__, and so on. It's __lightweight and easy__, unlike e-mails that have way too much overhead.

## How to run
### Hosted
One option is to simply use the hosted instance running at [https://apps.muetsch.io/webhook2telegram](https://apps.muetsch.io/webhook2telegram).  It only allows for a maximum of 240 requests per recipient per day. 

### Self-hosted
If you want to set this up on your own, do the following. You can either run the bot in long-polling- or webhook mode. For production use the latter option is recommended for [various reasons](https://core.telegram.org/bots/webhooks). However, you'll need a server with a static IP and s (self-signed) SSL certificate.

#### Compile from source
1. `git clone github.com/muety/webhook2telegram`
1. `GO111MODULE=on go build`
1. Long-polling mode: `./webhook2telegram -token <YOUR_BOTFATHER_TOKEN>`
1. Webhook mode: `./webhook2telegram -token <YOUR_BOTFATHER_TOKEN> -mode webhook`

#### Using Docker
```bash
$ docker volume create webhook2telegram_data
$ docker run -d -p 8080:8080 \
    -v webhook2telegram_data:/srv/data \
    -e "APP_TOKEN=<YOUR_BOTFATHER_TOKEN>" \
    -e "APP_MODE=webhook" \
    --name webhook2telegram \
    n1try/webhook2telegram
```

💡 It is recommended to either use `-useHttps` or set up a __reverse proxy__ like _nginx_ or [_Caddy_](https://caddyserver.com) to handle encryption.

### Additional parameters
* `-address` (`string`) – Network address (IPv4) to bind to. Defaults to `127.0.0.1`.
* `-address6` (`string`) – Network address (IPv6) to bind to. Defaults to `::1`.
* `-disableIPv6` (`bool`) – Whether to disable listening on both IPv4 and IPv6 interfaces. Defaults to `false`.
* `-port` (`int`) – TCP port to listen on. Defaults to `8080`.
* `-proxy` (`string`) – Proxy connection string to be used for long-polling mode. Defaults to none.
* `-useHttps` (`bool`) – Whether to use HTTPS. Defaults to `false`.
* `-certPath` (`string`) – Path of your SSL certificate when using webhook mode with `useHttp`. Default to none.
* `-keyPath` (`string`) – Path of your private SSL key when using webhook mode with `useHttp`. Default to none.
* `-dataDir` (`string`) – File system location where to store persistent data. Defaults to `.`.
* `-blacklist` (`string`) – Path to a user id blacklist file. Defaults to `blacklist.txt`.
* `-rateLimit` (`int`) – Maximum number of messages to be delivered to each recipient per hour. Defaults to `100`.
* `-metrics` (`bool`) – Whether to expose [Prometheus](https://prometheus.io) metrics under `/metrics`. Defaults to `false`.

## How to use
1. You need to get a token from the bot. Send a message with `/start` to the [Webhook2Telegram Bot](https://telegram.me/MiddleManBot) therefore.
2. Now you can use that token to make HTTP `POST` requests to `http://localhost:8080/api/messages` (replace localhost by the hostname of your server running the bot or mine as shown above) with a body that looks like this.

```
{
	"recipient_token": "3edf633a-eab0-45ea-9721-16c07bb8f245",
	"text": "*Hello World!* (yes, this is Markdown)",
	"type": "TEXT",
	"origin": "My lonely server script"
}
```

**NOTE:** If the field *type* is omitted then the `TEXT` type will be used as default, though this is not recommended as this may change in future versions.

More details can be found [here](/inlets).

### Inlets
Inlets provide a mechanism to pre-process incoming data that comes in a format different from what is normally expected by the bot. 

This is especially useful if data is sent by external, third-party applications which you cannot modify.

For instance, you might want to deliver alerts from [Prometheus' Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) as Telegram notifications. However, Alertmanager's [webhook requests](https://prometheus.io/docs/alerting/configuration/#webhook_config) look much different from WH2TG's default input format. To still make them fit, you can write an [Inlet](/inlets) to massage the data accordingly.

To directly address an inlet, request `http://localhost:8080/api/inlets/<inlet_name>`. Note that `/api/inlets/default` is equivalent to `/api/messages`.

Following inlets are currently available:

| Name         | Description                                                                                                 | Status |
|--------------|-------------------------------------------------------------------------------------------------------------|--------|
| `default`      | Simply passes the request through without any changes                                                       | ✅      |
| `alertmanager_webhook` | Consumes [Alertmanager webhook requests](https://prometheus.io/docs/alerting/configuration/#webhook_config) | ✅      |
| `bitbucket_webhook` | Accepts [Bitbucket webhook requests](https://confluence.atlassian.com/bitbucket/tutorial-create-and-trigger-a-webhook-747606432.html) to notify about a pipeline status change | ⏳      |
| `webmentionio_webhook` | Accepts [Webmention.io](https://webmention.io/) webhook requests to notify about a new Webmention of one of your articles | ✅      |

Further documentation about the individual inlets is available [here](/inlets).

### Metrics
Fundamental [Prometheus](https://prometheus) metrics are exposed under `/metrics`, if the `-metrics` flag gets passed. They include:
* `webhook2telegram_messages_total{origin="string", type="string"}` 
* `webhook2telegram_requests_total{success="string"}` 

## License
MIT @ [Ferdinand Mütsch](https://muetsch.io)
