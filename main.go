package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

//const BOT_API_TOKEN = "439545547:AAEiymDgV-ahktWxWs8dgTazanCmbenOTSg"
//const BASE_URL = "https://api.telegram.org/bot"
const BOT_API_TOKEN = ""
const BASE_URL = "https://requestb.in/19mgpmo1"
const API_URL = BASE_URL + BOT_API_TOKEN
const STORE_FILE = "store.json"

func sendMessage(recipientId, text string, isMarkdown bool) error {
	parseMode := ""
	if isMarkdown {
		parseMode = "Markdown"
	}
	m, err := json.Marshal(&TelegramOutMessage{ChatId: recipientId, Text: text, ParseMode: parseMode})
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(m))
	_, err = http.Post(API_URL, "application/json", reader)
	if err != nil {
		return err
	}
	return nil
}

func resolveToken(token string) string {
	return "1234"
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
		return
	}

	recipientId := resolveToken(m.RecipientToken)
	if len(recipientId) == 0 || len(m.Text) == 0 {
		w.WriteHeader(400)
		return
	}

	err = sendMessage(recipientId, m.Text, m.IsMarkdown)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

// TODO: get updates
// https://stackoverflow.com/questions/10152478/how-to-make-a-long-connection-with-http-client

func main() {
	ReadStoreFromJSON(STORE_FILE)

	// Exit handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	go func() {
		for _ = range c {
			fmt.Println("Flushing store.")
			FlushStoreToJSON(STORE_FILE)
			os.Exit(0)
		}
	}()

	http.HandleFunc("/api/messages", messageHandler)
	http.ListenAndServe(":8080", nil)
}
