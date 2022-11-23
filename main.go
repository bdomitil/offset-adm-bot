package main

import (
	"log"
	"offset-adm-bot/bitrix"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	upTimeout int = 60
	repList       = reportList{store: map[int64]report{}}
	Users         = UserList{mutex: new(sync.Mutex), store: map[int64]*user{}}
	BitrixU   bitrix.Profile
)

func SetUpTimeout(t int) {
	upTimeout = t
}

func (b *syncBot) Init() (updates tgbotapi.UpdatesChannel) {
	log.Println("Bot started")
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TTOKEN"))
	if err != nil {
		panic("Unable to start Telegram Bot, check if TTOKEN is available")
	}
	BitrixU, err = bitrix.Init(os.Getenv("BITRIX_TOKEN"))
	if err != nil {
		panic("Unable to connect Bitrix Api : " + err.Error())
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = upTimeout
	updates = bot.GetUpdatesChan(u)
	b.BotAPI = bot
	return updates
}

func main() {
	log.SetFlags(log.Ldate | log.Lmsgprefix | log.Lshortfile)
	bot := newSyncBot()
	updates := bot.Init()
	go reportsManager()
	go updateUserList(bot.Self.ID)
	go updateChats(bot, bot.Self.ID)
	for update := range updates {
		if update.Message == nil &&
			update.CallbackQuery == nil {
			continue
		}
		switch update.FromChat().Type {
		case "group":
			fallthrough
		case "supergroup":
			inGroupChat(bot, update)
		case "private":
			inPrivateChat(bot, update)
		default:
			log.Printf("message from : %s\nwith name %s\n", update.FromChat().Type,
				update.FromChat().Title)
		}
	}
}
