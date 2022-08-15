package main

import (
	"encoding/json"
	"fmt"
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

//Types:
//	offsetID
//	smile
type callbackJSON struct {
	Type string `json:"info"`
	Info string `json:"text"`
}

func isOffsetChat(title string) bool {
	status, _ := regexp.MatchString(`OF(\d{3}-\d{1,2})|OF(\d{3})`, title)
	return status
}

func getOffsets(title string) []string {
	regex, _ := regexp.Compile(`OF(\d{3}-\d{1,2})|OF(\d{3})`)
	return regex.FindAllString(title, -1)
}

func genReplyForMsgKeyboard(buttons ...string) []tgbotapi.KeyboardButton {
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

func select_OffId_Inline_keyboard(offs []string) (keyboard tgbotapi.InlineKeyboardMarkup) {
	var buttons []tgbotapi.InlineKeyboardButton
	for _, k := range offs {
		Callback := callbackJSON{Type: "offsetID", Info: k}
		call, _ := json.Marshal(Callback)
		button := tgbotapi.NewInlineKeyboardButtonData(k, string(call))
		buttons = append(buttons, button)
	}
	keyboard = tgbotapi.NewInlineKeyboardMarkup(buttons)
	return keyboard
}

//status = 2   : Reply is created and opened, needs info on problem
//status = 3   : Reply is ready to be published
//status = 4   : Needs to choose offset id
//status = 5   : offset id has been choosen
//status = 101 : Error openning new request, request is almost opened
//status = 102 : Error closing request, no opened requests found
//status = 200 : Reply is successfully closed
//status = 255 : Skip checking request status
func genReplyForMsg(update *tgbotapi.Update, status uint8) (reply tgbotapi.MessageConfig) {
	rep, ok := repList.findReport(update.FromChat().ID)
	if ok && status == 255 {
		status = rep.description.status
	}
	reply.ChatID = update.FromChat().ID
	// fmt.Printf("statusMSG = %d\n", status) //DEBUG
	switch status {
	case 2:
		reply.Text = replies["get_info_msg"] //Создание новой заявки
		reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyForMsgKeyboard(reportButtons["open"], reportButtons["close"]))
		rep.description.status = 3
		repList.putReport(update.FromChat().ID, rep)

	case 3:
		reply.Text = replies["request_filled_msg"] //Успешное заполнение заявки
		reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyForMsgKeyboard(reportButtons["open"], reportButtons["close"]))
		rep.description.comments = fmt.Sprintf("\tНомер аппарата - %s\n\n\tЖалоба:\n%s\n",
			rep.description.offID[0], update.Message.Text)
		if err := createTask(&BitrixU, rep.description); err != nil {
			reply.Text = "Ой((   Что-то пошло не так\nЯ уже передал сообщения администраторам\nМожете попробовать еще раз"
		}
		repList.close(update.Message.Chat.ID) // TODO: make it close by bitrix api
	case 4:
		reply.Text = "Пожалуйста выберите номер неисправного аппарата" //TODO change [terminal id] to [terminal location]
		reply.ReplyMarkup = select_OffId_Inline_keyboard(rep.description.offID)
		rep.description.status = 5
		repList.putReport(update.FromChat().ID, rep)
	case 101:
		reply.Text = "Пожалуйста завершите предыдущую заявку или нажмите + " + reportButtons["close"]
		reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyForMsgKeyboard(reportButtons["open"], reportButtons["close"]))
	case 102:
		reply.Text = "Нет открытых заявок\n" + reportButtons["open"] + " - чтоб открыть новую заявку"
		reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyForMsgKeyboard(reportButtons["open"], reportButtons["close"]))
	case 200:
		reply.Text = "Заявка успешко закрыта!"
		reply.ReplyToMessageID = repList.getReport(update.Message.Chat.ID).openMsgID
		reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(genReplyForMsgKeyboard(reportButtons["open"], reportButtons["close"]))
		repList.close(update.Message.Chat.ID)
	default:
		reply.Text = fmt.Sprintf("Ошибка при попытке генерации ответа, неизвестный статус заявки %s  %d!", getSmile("fail"), status)
	}
	return reply
}

func genReplyForCallback(update *tgbotapi.Update, status uint8, bot *tgbotapi.BotAPI) (reply tgbotapi.MessageConfig) {
	rep, ok := repList.findReport(update.FromChat().ID)
	if ok && status == 255 {
		status = rep.description.status
	}
	reply.ChatID = update.FromChat().ID
	// fmt.Printf("statusCallback = %d\n", status) //DEBUG
	switch status {
	case 5:
		rep.description.offID = []string{}
		rep.description.offID = append(rep.description.offID, update.CallbackQuery.Data)
		delmsg := tgbotapi.NewDeleteMessage(update.FromChat().ID, update.CallbackQuery.Message.MessageID)
		_, err := bot.Request(delmsg)
		if err != nil {
			log.Println(err.Error())
		}
		return genReplyForMsg(update, 2)
	default:
		errInline := tgbotapi.NewCallback(update.CallbackQuery.ID, fmt.Sprintf("Ошибка в статусе задача, статус = %d", status))
		bot.Request(errInline)
		reply.Text = getSmile("fail") + getSmile("fail")
	}
	return reply
}
