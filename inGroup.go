package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"offset-adm-bot/bitrix"
	"strconv"
)

func createTask(bit *bitrix.Profile, data reportForm) error {
	deals := bitrix.Get_deals()
	id, _ := bitrix.Get_deal_id_by_name(deals, data.offID)
	var newTask bitrix.Task
	newTask.Body.Deal_id = strconv.Itoa(id)
	newTask.Body.Description = data.comments
	newTask.Body.Responible_id = bit.Body.Id
	newTask.Body.Title = "Задача созданная " + bit.Body.Name
	err := bitrix.Task_add(&newTask)
	return err
}

func manageGroupChat(update *tgbotapi.Update) (reply tgbotapi.MessageConfig, err error) {

	if hasOpenReport(update) && openReports[update.Message.Chat.ID].creator != update.Message.From.ID { //return and not allow to any other reports ultil previous deletes
		return
	}

	reply = tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Text {
	case "/report":
		if hasOpenReport(update) {
			reply.Text = "Пожалуйста завершите предыдущую заявку или нажмите /close_report"
		} else {
			openReports[update.Message.Chat.ID] = openReport(update)
			reply.Text = genReply(update)
			reply.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(genReplyKeyboard("/close_report"))
		}
	case "/close_report":
		if hasOpenReport(update) {
			reply.Text = "Заявка успешко закрыта!"
			reply.ReplyToMessageID = openReports[update.Message.Chat.ID].openMsgID
			reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyKeyboard("/report"))
			delete(openReports, update.Message.Chat.ID)
		} else {
			reply.Text = "Нет открытых заявок, /report чтоб открыть новую заявку"
			reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyKeyboard("/report"))
		}
	default:
		if hasOpenReport(update) {
			fillReport(update)
			reply.Text = genReply(update)
			if openReports[update.Message.Chat.ID].isFilled {
				reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyKeyboard("/report"))
				if err = createTask(&BitrixU, openReports[update.Message.Chat.ID].description); err != nil {
					reply.Text = "Ой((   Что-то пошло не так\nЯ уже передал сообщения администраторам\nМожете попробовать еще раз"
				}
			}
		} else {
			reply.ChatID = 0
		}
	}

	return reply, err
}
