package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type reportList struct {
	store map[int64]report
}

type reportForm struct {
	offID    []string
	comments string
	status   uint8
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

var replies = map[string]string{
	"hello_msg": "Здраствуйте, меня зовут Оффсетик, я бот техподдержки компании OFFSET\n\n", //Приветсвенное сообщение
	"get_info_msg": `Cейчас разберемся что у вас случилось 
	Опишите пожалуйста проблему максимально подробно, чем больше деталей тем лучше!
	
	При описании проблемы ответьте на следующие вопросы 
	
	1) Когда возникла проблема и в следствии чего ? 
	2) Какие действия были предприняты ? 
	3) По возможности направьте файл - при распечатке которого возникла проблема`, //Сообщение о сборе информации

	"request_filled_msg": "Ваша заявка успешно принята, начинаем решать вопрос", //Сообщение о принятии заявки
}

func (rep *reportList) getReport(id int64) (r report) {
	r = rep.store[id]
	return
}

func (rep *reportList) findReport(id int64) (r report, ok bool) {
	r, ok = rep.store[id]
	return
}

func (rep *reportList) getStore() (store *map[int64]report) {
	return &rep.store
}

func (rep *reportList) putReport(id int64, r report) {
	rep.store[id] = r
}

func (rep *reportList) isOpen(update *tgbotapi.Update) bool {
	ok := false
	_, ok = rep.store[update.FromChat().ID]
	return ok
}

func (rep *reportList) close(id int64) {
	if _, ok := rep.store[id]; ok {
		delete(*rep.getStore(), id)
	}
}

func newReport(update *tgbotapi.Update) (newR report) {
	newR = report{}
	newR.isFilled = false
	newR.creator = update.SentFrom().ID
	newR.creation_time = time.Now()
	newR.channel_name = update.FromChat().Title
	newR.channel_id = update.FromChat().ID
	newR.openMsgID = update.Message.MessageID
	newR.description = reportForm{}
	newR.description.offID = getOffsets(newR.channel_name)
	if len(newR.description.offID) > 1 {
		newR.description.status = 4
	} else {
		newR.description.status = 2
	}
	return newR
}
