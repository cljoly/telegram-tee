package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

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

// Parse chat ID from command line arguments
func parseChatID(args []string) (chatIDs []int, err error) {
	chatIDs = make([]int, len(args))
	err = nil
	var id int
	for i := 0; i < len(args); i++ {
		id, err = strconv.Atoi(args[i])
		chatIDs[i] = id
		if err != nil {
			return
		}
	}
	return
}

// Message handler replying by giving a set of chat ID
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
func messageWriter(bot *tgbotapi.BotAPI, chatIDs []int) {
	var msgTxt string
	var msg tgbotapi.MessageConfig
	scanner := bufio.NewScanner(os.Stdin)
	log.Printf("Scannig…")
	for scanner.Scan() {
		msgTxt = scanner.Text()
		log.Printf("Scanned '%s'", msgTxt)
		for chatID := range chatIDs {
			msg = tgbotapi.NewMessage(int64(chatID), msgTxt)

			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message from stdin: '%s'", err)
			}
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

	go updateLoop(bot, replyChatID)

	if len(os.Args) < 2 {
		fmt.Println("I’m listening, I will answer to everyone with the chatID of our chat. You can then make me write the content of stdin to a chat, by giving me its chatID as argument.")
	} else {
		ids, _ := parseChatID(os.Args[2:])
		messageWriter(bot, ids)
	}
}
