package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	botToken := os.Getenv("TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("creating bot API: %v", err) // You should add better error handling than this!
	}

	bot.Debug = true // Has the library display every request and response.
	config := tgbotapi.NewWebhook("https://hpmorcitation.appspot.com/" + botToken)
	resp, err := bot.SetWebhook(config)
	if err != nil {
		log.Fatalf("setting webhook: %v", err)
	}
	if !resp.Ok {
		log.Fatalf("setting webhook: %s", resp.Description)
	}
}
