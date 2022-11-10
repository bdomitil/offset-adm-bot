package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//TODO change cacheGroup url
func cacheGroup(Chat *tgbotapi.Chat, selfId int64) error {
	newChat := newChat(*Chat, selfId)
	err := saveChat(newChat)
	return err
}
