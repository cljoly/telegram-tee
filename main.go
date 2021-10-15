package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Create a bot with token
func login(token string) (bot *tgbotapi.BotAPI) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic("login error")
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
		fmt.Fprintln(os.Stderr, err)
		panic("Error getting updates")
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

	fmt.Fprintf(os.Stderr, "[%s(%s)] %s", update.Message.From.UserName, chatID, update.Message.Text)

	msgText := fmt.Sprint("chatID: ", chatID)

	msg := tgbotapi.NewMessage(chatID, msgText)
	// msg.ReplyToMessageID = update.Message.MessageID

	_, err := bot.Send(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error sending message: ", msgText)
		panic("error replying")
	}
}

// Create a message handler to write content of stdout (each line) to someone
func messageWriter(bot *tgbotapi.BotAPI, chatIDs []int) {
	var msgTxt string
	var msg tgbotapi.MessageConfig
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprintf(os.Stderr, "Scannig…")
	for scanner.Scan() {
		msgTxt = scanner.Text()
		fmt.Fprintf(os.Stderr, "Scanned '%s'", msgTxt)
		for _, chatID := range chatIDs {
			msg = tgbotapi.NewMessage(int64(chatID), msgTxt)
			_, err := bot.Send(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error sending message to %d from stdin: '%s'", chatID, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return
}

func main() {
	token := os.Getenv("TLGCLI_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "You need to set TLGCLI_TOKEN in your env")
		return
	}
	bot := login(token)
	bot.Debug = false
	fmt.Fprintf(os.Stderr, "Authorized on account %s", bot.Self.UserName)

	go updateLoop(bot, replyChatID)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "I’m listening, I will answer to everyone with the chatID of our chat. You can then make me write the content of stdin to a chat, by giving me its chatID as argument.")
	} else {
		ids, _ := parseChatID(os.Args[1:])
		messageWriter(bot, ids)
	}
}
