package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var replies = [...]string{
	"",
	"Здраствуйте, меня зовут Оффсетик, я бот техподдержки компании OFFSET, сейчас разберемся что у вас случилось\n\n" +
		"Опишите пожалуйста проблему польность, чем больше деталей тем лучше!",
	"Напишите были ли предприняты какие-то действия самостоятельно, если да то, что к чему это привело",
	"Ваша заявка успешно принята, начинаем решать вопрос",
}

type reportForm struct {
	offID        string
	comments     string
	moves_to_fix string
	size         uint8
	//TODO: add field for media data
}

type report struct {
	creator       int64
	description   reportForm
	channel_id    int64
	channel_name  string
	creation_time time.Time
	isFilled      bool
	openMsgID     int
}

func fillReport(update *tgbotapi.Update) {

	var rep report = openReports[update.Message.Chat.ID]
	switch rep.description.size {
	case 1:
		rep.description.comments = update.Message.Text
	case 2:
		rep.description.moves_to_fix = update.Message.Text
		fallthrough
		//TODO: case 3 : rep.description.media_facts = somemedia
	case 3:
		rep.isFilled = true
	}
	rep.description.size++
	openReports[update.Message.Chat.ID] = rep
}

func genReply(update *tgbotapi.Update) (reply string) {
	var rep report = openReports[update.Message.Chat.ID]

	switch rep.description.size {
	case 1:
		reply = replies[1]
	case 2:
		reply = replies[2]
	case 3:
		reply = replies[3]
	default:
		reply = "Ошибка при попытке генерации ответа"
	}
	return reply
}

func openReport(update *tgbotapi.Update) (newR report) {
	newR.isFilled = false
	newR.creator = update.Message.From.ID
	newR.creation_time = time.Now()
	newR.channel_name = update.Message.Chat.Title
	newR.channel_id = update.Message.Chat.ID
	newR.openMsgID = update.Message.MessageID
	newR.description = reportForm{}
	newR.description.size = 1
	_, OffId := isOffsetChat(newR.channel_name)
	newR.description.offID = OffId
	return newR
}
