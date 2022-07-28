package main

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var openReports = map[int64]report{}

func main() {

	fmt.Println("Bot started")
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TTOKEN"))
	if err != nil {
		panic("Unable to start Telegram Bot, check if TTOKEN is available")
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {

		if update.Message == nil {
			continue
		}

		var newMsg tgbotapi.MessageConfig
		if isChat, _ := isOffsetChat(update.Message.Chat.Title); update.Message.Chat.IsGroup() && isChat {
			// newMsg.Text = fmt.Sprintf("sender name = %s, sender id = %d, message chat name  = %s, message chat id = %d ",
			// 	update.Message.From.FirstName, update.Message.From.ID, update.Message.Chat.Title, update.Message.Chat.ID)
			// bot.Send(newMsg)
			newMsg, err = manageGroupChat(&update)
		} else if update.Message.Chat.IsPrivate() {
			newMsg.Text = "хочешь приват?"
			newMsg.ChatID = update.Message.Chat.ID
		}
		bot.Send(newMsg)
	}
}

// func main() {

// 	api_test()
// }
