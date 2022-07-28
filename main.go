package main

import (
	"fmt"
	"log"
	"offset-adm-bot/bitrix"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var openReports = map[int64]report{}
var BitrixU bitrix.Profile

func main() {
	fmt.Println("Bot started")
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TTOKEN"))
	if err != nil {
		panic("Unable to start Telegram Bot, check if TTOKEN is available")
	}
	bitrixU, err := bitrix.Init(os.Getenv("bitrix_api_url"))
	if err != nil {
		panic("Unable to connect Bitrix Api : " + err.Error())
	}
	BitrixU = bitrixU
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {

		if update.Message == nil {
			continue
		}
		var newMsg tgbotapi.MessageConfig
		if isChat, _ := isOffsetChat(update.Message.Chat.Title); update.Message.Chat.IsGroup() && isChat {
			newMsg, err = manageGroupChat(&update)
		} else if update.Message.Chat.IsPrivate() {
			newMsg.Text = "хочешь приват?"
			newMsg.ChatID = update.Message.Chat.ID
		}
		if err != nil {
			log.Println(err.Error())
		}
		bot.Send(newMsg)
	}
}

// func main() {

// 	api_test()
// }
