package main

import (
	// "flag"
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Create a bot with token
func login(token string) (bot *tgbotapi.BotAPI) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return
}

// Type of functions handling updates
type handler func(*tgbotapi.BotAPI, tgbotapi.Update)

// Get updates continuously for a bot and pass it to a handler
func updateLoop(bot *tgbotapi.BotAPI, h handler) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic("Error getting updates")
	}

	for update := range updates {
		h(bot, update)
	}
}

// Message handler replying by giving current chat ID
func replyChatID(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	chatID := update.Message.Chat.ID

	log.Printf("[%s(%s)] %s", update.Message.From.UserName, chatID, update.Message.Text)

	msgText := fmt.Sprint("chatID: ", chatID)

	msg := tgbotapi.NewMessage(chatID, msgText)
	// msg.ReplyToMessageID = update.Message.MessageID

	_, err := bot.Send(msg)
	if err != nil {
		log.Panic("Error sending message: ", msgText)
	}
}

func main() {
	bot := login(os.Getenv("TLGCLI_TOKEN"))

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateLoop(bot, replyChatID)
}
