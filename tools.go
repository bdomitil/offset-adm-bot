package main

import (
	"log"
	"os"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func hasOpenReport(update *tgbotapi.Update) bool {
	report, ok := openReports[update.Message.Chat.ID]
	if ok && !report.isFilled {
		return true
	}
	return false
}

func isOffsetChat(chatName string) (bool, string) {
	regex, _ := regexp.Compile(`^OF\d\d\d.`)
	status, _ := regexp.MatchString(regex.String(), chatName)
	return status, regex.FindString(chatName)
}

func genReplyKeyboard(buttons ...string) []tgbotapi.KeyboardButton {
	keyboards := make([]tgbotapi.KeyboardButton, 1, 2)
	for _, button := range buttons {
		keyboards = append(keyboards, tgbotapi.NewKeyboardButton(button))
	}
	return keyboards
}

func sendAdminErroMsg(bot *tgbotapi.BotAPI, text string) {
	admin_id, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	if err != nil || admin_id == 0 {
		log.Fatalf("Admin telegram chat id is false")
	}
	var newMsg tgbotapi.MessageConfig
	newMsg.ChatID = int64(admin_id)
	newMsg.Text = text
	bot.Send(newMsg)
}
