package main

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"offset-adm-bot/bitrix"
	"strconv"
)

var reportButtons = map[string]string{
	"open":   "openreport",
	"delete": "Удалить обращение",
	"finish": "Завершить обращение",
	"close":  "closereport",
	"start":  "start"}

func createTask(bit *bitrix.Profile, data reportForm) error {
	deals := bitrix.Get_deals()
	id, _ := bitrix.Get_deal_id_by_name(deals, data.offID[0]) //TODO: fix it
	var newTask bitrix.Task
	newTask.Fields.Deal_id = fmt.Sprintf("D_%s", strconv.Itoa(id))
	newTask.Fields.Description = data.comments
	newTask.Fields.Responible_id = bit.Result.Id
	newTask.Fields.Title = fmt.Sprintf("Задача созданная  %s %s", bit.Result.Name, bit.Result.Last_name)
	err := bitrix.Task_add(&newTask)
	return err
}

func manageGroupChat(update *tgbotapi.Update, bot *tgbotapi.BotAPI) (reply tgbotapi.MessageConfig, err error) {

	if (repList.isOpen(update) && repList.getReport(update.FromChat().ID).creator !=
		update.SentFrom().ID) || update.SentFrom().IsBot { //return and not allow to any other reports ultil previous deletes
		return
	}
	reply = tgbotapi.NewMessage(update.FromChat().ID, "")
	if update.Message != nil { //Client sent message
		switch update.Message.Command() {
		case reportButtons["open"]:
			if repList.isOpen(update) {
				reply = genReplyForMsg(update, 101)
			} else {
				repList.putReport(update.Message.Chat.ID, newReport(update)) //creating new report
				reply = genReplyForMsg(update, 255)
			}
		case reportButtons["close"]:
			if repList.isOpen(update) {
				reply = genReplyForMsg(update, 200)
			} else {
				reply = genReplyForMsg(update, 102)
			}
		case reportButtons["start"]:
			reply = genReplyForMsg(update, 1)
		default:
			{
				if repList.isOpen(update) {
					reply = genReplyForMsg(update, 255)
				}
			}
		}
	} else if update.CallbackQuery != nil { //Client sent callback
		var callback callbackJSON
		err := json.Unmarshal([]byte(update.CallbackData()), &callback)
		if err != nil {
			return reply, fmt.Errorf("unknown callback : %s", update.CallbackData())
		}
		update.CallbackQuery.Data = callback.Info
		switch callback.Type {
		case "offsetID":
			{
				reply = genReplyForCallback(update, 255, bot)
			}
		}
	}
	return reply, err
}
