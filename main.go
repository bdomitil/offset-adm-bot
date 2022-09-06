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
		if update.Message == nil &&
			update.CallbackQuery == nil {
			continue
		}
		var newMsg tgbotapi.MessageConfig
		isChat := isOffsetChat(update.FromChat().Title)
		if (update.FromChat().IsGroup() ||
			update.FromChat().IsSuperGroup()) &&
			isChat { //allows just offset groups
			if update.CallbackQuery != nil {
				newMsg, err = manageGroupChat(&update, bot) //manage callback queries commands
			} else if len(update.Message.NewChatMembers) > 0 { //manage new chat members
				newMsg, _ = manageUserEntry(bot, &update)
			} else if update.Message != nil && len(update.Message.Text) > 0 { //manage text messages commands
				newMsg, err = manageGroupChat(&update, bot)
			}
			if err != nil && err.Error() == "skip" {
				continue
			} else if err != nil {
				sendAdminErroMsg(bot, err.Error())
				continue
			}
		} else if update.FromChat().IsPrivate() { //manage all private chats
			newMsg.Text = "Я пока еще не умею общаться так, но очень скоро научусь! дождись меня"
			newMsg.ChatID = update.FromChat().ID
		}
		if newMsg.Text != "" {
			_, err = bot.Send(newMsg)
		}
		if err != nil {
			log.Println(err.Error())
		}
	}
}
