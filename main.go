package main

import (
	"fmt"
	"log"
	"offset-adm-bot/bitrix"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	upTimeout int = 60
	repList       = reportList{store: map[int64]report{}}
	Users         = map[int64]user{}
	BitrixU   bitrix.Profile
	// bot       *tgbotapi.BotAPI
	// update    *tgbotapi.Update
)

func SetUpTimeout(t int) {
	upTimeout = t
}

func Init() (updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	fmt.Println("Bot started")
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TTOKEN"))
	if err != nil {
		panic("Unable to start Telegram Bot, check if TTOKEN is available")
	}
	BitrixU, err := bitrix.Init(os.Getenv("BITRIX_TOKEN"))
	if err != nil {
		panic("Unable to connect Bitrix Api : " + err.Error())
	}
	_ = BitrixU
	u := tgbotapi.NewUpdate(0)
	u.Timeout = upTimeout
	updates = bot.GetUpdatesChan(u)
	return updates, bot
}

func main() {
	var err error
	updates, bot := Init()
	for update := range updates {
		if update.Message == nil &&
			update.CallbackQuery == nil {
			continue
		}
		updateUserList(bot.Self.ID)
		var newMsg tgbotapi.MessageConfig
		newMsg.ChatID = update.FromChat().ID
		isChat := isOffsetChat(update.FromChat().Title)
		switch {
		case (update.FromChat().IsGroup() || update.FromChat().IsSuperGroup()) &&
			isChat:
			{ //allows just offset groups
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
				bot.Send(newMsg)
			}
		case update.FromChat().IsPrivate() && isUserAuthed(update.FromChat().ID): //manage all private chats
			user := Users[update.FromChat().ID]
			user.adminPanelExec(bot, &update)
			Users[update.FromChat().ID] = user
		default:
			newMsg.Text = "Я пока еще не умею общаться так, но очень скоро научусь! дождись меня"
			bot.Send(newMsg)
		}
		if err != nil {
			log.Println(err.Error())
		}
	}
}
