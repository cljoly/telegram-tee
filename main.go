package main

import (
	"bufio"
	"flag"
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

// Create a message handler to write content of stdout (each line) to someone
func messageWriter(bot *tgbotapi.BotAPI, chatID int64) {
	var msgTxt string
	var msg tgbotapi.MessageConfig
	scanner := bufio.NewScanner(os.Stdin)
	log.Printf("Scannig…")
	for scanner.Scan() {
		msgTxt = scanner.Text()
		log.Printf("Scanned '%s'", msgTxt)
		msg = tgbotapi.NewMessage(chatID, msgTxt)

		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message from stdin: '%s'", err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return
}

func main() {
	bot := login(os.Getenv("TLGCLI_TOKEN"))
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	chatID := flag.NewFlagSet("chatid", flag.ExitOnError)

	write := flag.NewFlagSet("write", flag.ExitOnError)
	var id int64
	write.Int64Var(&id, "to", -1, "Id of the chat to write to")

	if len(os.Args) < 2 {
		fmt.Println("TODO Write doc")
		os.Exit(2)
	}
	var h handler = nil
	switch os.Args[1] {
	case "chatid":
		chatID.Parse(os.Args[2:])
		h = replyChatID
	case "write":
		write.Parse(os.Args[2:])
		messageWriter(bot, id)
	}
	if h != nil {
		updateLoop(bot, h)
	}
}
