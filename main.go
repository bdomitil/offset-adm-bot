package main

import (
	"fmt"
	"log"
	"offset-adm-bot/bitrix"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var repList = reportList{store: map[int64]report{}}
var BitrixU bitrix.Profile

func main() {

	fmt.Println("Bot started")
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TTOKEN"))
	if err != nil {
		panic("Unable to start Telegram Bot, check if TTOKEN is available")
	}
	bitrixU, err := bitrix.Init(os.Getenv("BITRIX_TOKEN"))
	if err != nil {
		panic("Unable to connect Bitrix Api : " + err.Error())
	}
	BitrixU = bitrixU
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil ||
			(update.Message.Text == "" && len(update.Message.NewChatMembers) == 0) &&
				update.CallbackQuery == nil {
			tgbotapi.NewBotCommandScopeAllGroupChats()
			continue
		}
		var newMsg tgbotapi.MessageConfig

		if isChat := isOffsetChat(update.FromChat().Title); (update.FromChat().IsGroup() || update.FromChat().IsSuperGroup()) && isChat {
			if len(update.Message.NewChatMembers) > 0 && bot.Self.ID == update.Message.NewChatMembers[0].ID {
				update.Message.Text = "start"
			}
			newMsg, err = manageGroupChat(&update, bot)
			if err != nil && err.Error() == "skip" {
				continue
			} else if err != nil {
				sendAdminErroMsg(bot, err.Error())
				continue
			}
			if newMsg.Text != "" {
				_, err = bot.Send(newMsg)
			}
		} else if update.FromChat().IsPrivate() {
			newMsg.Text = "Я пока еще не умею общаться так, но очень скоро научусь! дождись меня"
			newMsg.ChatID = update.FromChat().ID
			_, err = bot.Send(newMsg)
		}
		if err != nil {
			log.Println(err.Error())
		}
	}
}
