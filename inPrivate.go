package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func isUserAuthed(id int64) bool {
	Users.mutex.Lock()
	_, ok := Users.store[id]
	Users.mutex.Unlock()
	return ok
}

func inPrivateChat(bot *syncBot, update tgbotapi.Update) {
	if isUserAuthed(update.FromChat().ID) { //manage all private chats
		go func() {
			user, _ := Users.PopUser(update.FromChat().ID)
			user.adminPanelExec(bot, update)
			Users.PushUser(user)
		}()
	}
}

func (u *user) adminPanelExec(bot *syncBot, update tgbotapi.Update) {
	u.Block()
	cmd, err := u.newCmd(Button(update.Message.Text))
	if err != nil {
		log.Println(err)
		return
	}
	u.cmd = cmd
	cmd.exec(bot, update)
	if u.cmd.getState() == Closed {
		u.prevCmd = nil
	} else {
		u.prevCmd = cmd.copy()
	}
	u.Unblock()
}
