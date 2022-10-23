package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func isUserAuthed(id int64) bool {
	_, ok := Users[id]
	return ok
}

func (u *user) adminPanelExec(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	cmd, err := u.newCmd(update.Message.Text)
	if err != nil {
		log.Println(err)
	}
	u.cmd = cmd
	cmd.exec(bot, update)
	u.prevCmd = cmd.copy()
}
