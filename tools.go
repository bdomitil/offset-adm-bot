package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var smiles = map[string][]byte{
	"laugh": []byte("\xF0\x9F\x98\x81"),
	"done":  []byte("\xE2\x9C\x85"),
	"fail":  []byte("\xE2\x9D\x8C"),
	"hand":  []byte("\xE2\x9C\x8B"),
	"comp":  []byte("\xF0\x9F\x93\x87"),
}

var initText string = `Доброго времени суток☀️🌙

Меня зовут бот Оффсетик 
Я добавлен в Ваш чат для решения всех технических проблем в дальнейшем🦸🏻‍♂️
Можете обращаться ко мне за поддержкой 24\7, я буду рад помочь 🤗
Надеюсь на нашу долгую плодотворную работу ✨`

var last_user_up time.Time = time.Now()

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

func sendAdminErroMsg(bot *syncBot, text string) {
	admin_id, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	if err != nil || admin_id == 0 {
		log.Fatalf("Admin telegram chat id is false")
	}
	var newMsg tgbotapi.MessageConfig
	newMsg.ChatID = int64(admin_id)
	newMsg.Text = text
	bot.syncSend(newMsg)
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
//status = 1   : init message of bot
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
	case 1:
		reply.Text = initText
	case 2:
		reply.Text = replies["get_info_msg"] //Создание новой заявки
		rep.description.status = 3
		repList.putReport(update.FromChat().ID, rep)

	case 3:
		reply.Text = replies["request_filled_msg"] //Успешное заполнение заявки
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
	case 102:
		reply.Text = "Нет открытых заявок\n" + reportButtons["open"] + " - чтоб открыть новую заявку"
	case 200:
		reply.Text = "Заявка успешко закрыта!"
		reply.ReplyToMessageID = repList.getReport(update.Message.Chat.ID).openMsgID
		repList.close(update.Message.Chat.ID)
	default:
		if status == 5 {
			reply.Text = fmt.Sprintln("Пожалуйста завершите выбор неисправного аппарата ") + getSmile("hand")
		} else {
			reply.Text = fmt.Sprintf("Ошибка при попытке генерации ответа, неизвестный статус заявки %s  %d!", getSmile("fail"), status)
		}
	}
	return reply
}

func NewResizeOneTimeReplyKeyboard(buttons ...string) (keyboard tgbotapi.ReplyKeyboardMarkup) {
	row := make([]tgbotapi.KeyboardButton, 0)
	rows := make([][]tgbotapi.KeyboardButton, 0)
	for i, b := range buttons {
		butt := tgbotapi.NewKeyboardButton(b)
		row = append(row, butt)
		if (i+1)%3 == 0 {
			rows = append(rows, row)
			row = make([]tgbotapi.KeyboardButton, 0)
		}
	}
	rows = append(rows, row)
	keyboard = tgbotapi.NewReplyKeyboard(rows...)
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true
	return
}
func genReplyForCallback(update *tgbotapi.Update, status uint8, bot *syncBot) (reply tgbotapi.MessageConfig) {
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

func updateUserList(botID int64) {
	currentT := time.Now()
	upPeriod := time.Minute * 5
	// url := fmt.Sprintf("http://localhost:3334/user/list/%d", botID)
	url := fmt.Sprintf("http://tg_cache:3334/user/list/%d", botID)

	if currentT.After(last_user_up) {
		last_user_up = time.Now().Add(upPeriod)
		var newUsers []user
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")
		response, err := http.DefaultClient.Do(req)

		if err != nil || response.StatusCode != 200 {
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
		defer response.Body.Close()
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err.Error())
			return
		}
		err = json.Unmarshal(responseBody, &newUsers)
		if err != nil {
			log.Println(err.Error())
			return
		}
		upUsers := make(map[int64]user)
		for _, u := range newUsers {
			if _, ok := Users[u.User_id]; ok {
				u.cmd = Users[u.User_id].cmd
				u.prevCmd = Users[u.User_id].prevCmd
			}
			upUsers[u.User_id] = u
		}
		for k := range Users {
			delete(Users, k)
		}
		Users = upUsers
	}
}

func getChatsForBot(botID int64) (chats []chat, err error) {
	resp, err := http.Get(fmt.Sprintf("http://tg_cache:3334/chat/list/%d", botID)) //TODO change to config parse
	// resp, err := http.Get(fmt.Sprintf("http://localhost:3334/chat/list/%d", botID)) //TODO change to config parse
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}
	err = json.Unmarshal(body, &chats)
	if err != nil {
		return
	}
	return chats, nil
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

func (bot *syncBot) syncSend(value tgbotapi.Chattable) (msg tgbotapi.Message, err error) {
	bot.mutex.Lock()
	msg, err = bot.Send(value)
	time.Sleep(time.Millisecond * 300)
	bot.mutex.Unlock()
	return
}

func newSyncBot() (bot *syncBot) {
	return new(syncBot)
}

func getDepartment(title string) (dep string) {
	deps := strings.Split(title, " ")
	deps1 := strings.Split(title, ".")

	if deps[len(deps)-1] == "ОС" {
		return "ОС"
	}
	if deps1[len(deps1)-1] == "ОС" {
		return "ОС"
	}
	return "ОЗ"
}

func reportsManager() {
	for {
		repList.mutex.Lock()
		for _, r := range repList.store {
			if time.Now().After(r.creation_time.Add(time.Minute * 30)) {
				delete(repList.store, r.chat_id)
			}
		}
		repList.mutex.Unlock()
		time.Sleep(time.Minute * 3)
	}
}
