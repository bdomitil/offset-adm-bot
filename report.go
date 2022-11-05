package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var replies = map[string]string{
	"hello_msg": "Здраствуйте, меня зовут Оффсетик, я бот техподдержки компании OFFSET\n\n", //Приветсвенное сообщение
	"get_info_msg": `Cейчас разберемся что у вас случилось 
	Опишите пожалуйста проблему максимально подробно, чем больше деталей тем лучше!
	
	При описании проблемы ответьте на следующие вопросы 
	
	1) Когда возникла проблема и в следствии чего ? 
	2) Какие действия были предприняты ? `, //Сообщение о сборе информации
	//3) По возможности направьте файл - при распечатке которого возникла проблемe
	"request_filled_msg": "Ваша заявка успешно принята, начинаем решать вопрос", //Сообщение о принятии заявки
}

func (rep *reportList) getReport(id int64) (r report) {
	rep.mutex.Lock()
	r = rep.store[id]
	rep.mutex.Unlock()
	return
}

func (rep *reportList) findReport(id int64) (r report, ok bool) {
	rep.mutex.Lock()
	r, ok = rep.store[id]
	rep.mutex.Unlock()
	return
}

func (rep *reportList) getStore() (store *map[int64]report) {
	rep.mutex.Lock()
	st := &rep.store
	rep.mutex.Unlock()
	return st
}

func (rep *reportList) putReport(id int64, r report) {
	rep.mutex.Lock()
	rep.store[id] = r
	rep.mutex.Unlock()
}

func (rep *reportList) isOpen(update *tgbotapi.Update) bool {
	rep.mutex.Lock()
	ok := false
	_, ok = rep.store[update.FromChat().ID]
	rep.mutex.Unlock()
	return ok
}

func (rep *reportList) close(id int64) {
	rep.mutex.Lock()
	if _, ok := rep.store[id]; ok {
		delete(*rep.getStore(), id)
	}
	rep.mutex.Unlock()
}

func newReport(update *tgbotapi.Update) (newR report) {
	newR = report{}
	newR.isFilled = false
	newR.creator = update.SentFrom().ID
	newR.creation_time = time.Now()
	newR.channel_name = update.FromChat().Title
	newR.chat_id = update.FromChat().ID
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
