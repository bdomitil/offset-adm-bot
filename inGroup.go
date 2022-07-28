package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

			// report := openReports[update.Message.Chat.ID]
			// for _, s := range openReports[update.Message.Chat.ID].description {
			// 	fmt.Printf("chat id = %d , desc = %s\n", report.channel_id, s)
			// }
			fillReport(update)
			reply.Text = genReply(update)
			if openReports[update.Message.Chat.ID].isFilled {
				reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyKeyboard("/report"))
				deals := get_deals()
				id, _ := get_deal_id_by_name(deals, openReports[update.Message.Chat.ID].description.offID)
				for _, dealx := range deals.Body {
					fmt.Printf("title = %s\nid = %s\n", dealx.Title, dealx.Id)
				}
				fmt.Printf("searched id = %d, name = %s\n", id, openReports[update.Message.Chat.ID].description.offID)
				data := openReports[update.Message.Chat.ID].description
				add_task_to_deal(id, &data)
			}
		} else {
			reply.ChatID = 0
		}
	}

	return reply, nil
}
