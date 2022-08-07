package main

import (
	"log"
	"os"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var smiles = map[string][]byte{
	"laugh": []byte("\xF0\x9F\x98\x81"),
	"done":  []byte("\xE2\x9C\x85"),
	"fail":  []byte("\xE2\x9D\x8C"),
	"hand":  []byte("\xE2\x9C\x8B"),
	"comp":  []byte("\xF0\x9F\x93\x87"),
}

func hasOpenReport(update *tgbotapi.Update) bool {
	report, ok := openReports[update.Message.Chat.ID]
	if ok && !report.isFilled {
		return true
	}
	return false
}

func isOffsetChat(title string) bool {
	status, _ := regexp.MatchString(`OF\d\d\d`, title)
	return status
}

func getOffsets(title string) []string {
	regex, _ := regexp.Compile(`OF\d\d\d`)
	return regex.FindAllString(title, -1)
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

func getSmile(s string) string {
	return string(smiles[s])
}
