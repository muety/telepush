# telegram-middleman-bot

![](http://i.imgur.com/lvshgaj.png)

__I'm the [@MiddleMan](https://telegram.me/MiddleManBot) bot! I sit in the middle between whatever you want to send you a message and your Telegram.__

I translate simple JSON HTTP requests into Telegram push messages that you will get on your Smartphone, PC or whatever Telegram client you have.

## Why might this be useful?
This is especially useful for __developers or sysadmins__. Imagine you want some kind of reporting from your application or server, like a daily report including some statistics. You don't want to actively look it up on a website but you want to receive it in a __passive fashion__. Just like getting an e-mail. But come on, let's be honest. __E-Mails are so 2010__. And they require your little server-side script to include some SMTP library and connect to a mail server. That's __too heavyweight__ just to __get some short information__. Personally, I have a Python script running on my server which gathers some statistics from log files and databases and regularly sends me a Telegram message.

If you develop those thoughts further, this could potentially __replace any kind of e-mail notifications__ - be it the message that someone has answered to your __forum post__, your favorite game is __now on sale at Steam__, and so on. It's __lightweight and easy__, unlike e-mails that have way too much overhead.

## How to run it?
You can either set up your own instance or use mine, which is running at [http://middleman.ferdinand-muetsch.de](http://middleman.ferdinand-muetsch.de). If you want to set this up on your own, do the following.

1. Make sure u have the latest version of Go installed.
2. `go get github.com/n1try/telegram-middleman-bot`
3. `cd <YOUR_GO_WORKSPACE_PATH>/src/github.com/n1try/telegram-middleman-bot`
4. `go get ./...`
5. Insert your `BOT_API_TOKEN`, which you got from the [@BotFather](https://telegram.me/BotFather) when registering your bot, in `main.go`
6. `go build .`
7. `./telegram-middleman-bot`

## How to use it?
1. You need to get a token from the bot. Send a message with `/start` to the [@MiddleManBot](https://telegram.me/MiddleManBot) therefore.
2. Now you can use that token to make HTTP POST requests to `http://localhost:8080/api/messages` (replace localhost by the hostname of your server running the bot or mine as shown above) with a body that looks like this.

```
{
	"recipient_token": "3edf633a-eab0-45ea-9721-16c07bb8f245",
	"text": "__Hello World!__ (yes, this is Markdown)",
	"origin": "My lonely server script"
}
```
## License
MIT @ [Ferdinand Mütsch](https://ferdinand-muetsch.de)