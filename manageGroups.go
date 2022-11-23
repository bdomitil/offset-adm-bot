package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//TODO change cacheGroup url
func cacheGroup(Chat *tgbotapi.Chat, selfId int64) error {
	newChat := newChat(*Chat, selfId)
	err := saveChat(newChat)
	return err
}

func inGroupChat(bot *syncBot, update tgbotapi.Update) {

	var err error
	if isOffsetChat(update.FromChat().Title) { //allows just offset groups
		newMsg := tgbotapi.NewMessage(update.FromChat().ID, "")
		if update.CallbackQuery != nil {
			newMsg, err = manageGroupChat(&update, bot) //manage callback queries commands
		} else if len(update.Message.NewChatMembers) > 0 { //manage new chat members
			newMsg, err = manageUserEntry(bot, &update)
		} else if update.Message != nil && len(update.Message.Text) > 0 { //manage text messages commands
			newMsg, err = manageGroupChat(&update, bot)
		}
		if err != nil && err.Error() == "skip" {
			return
		} else if err != nil {
			sendAdminErroMsg(bot, err.Error())
			log.Println(err)
			return
		}
		if newMsg.Text != "" {
			_, err = bot.syncSend(newMsg)
		}
		if err != nil {
			log.Println(err)
		}
	}
}
