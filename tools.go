package main

import (
	"bytes"
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
	"sync"
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
		log.Println("Admin telegram chat id is false")
		return
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

func NewResizeOneTimeReplyKeyboard(buttons []Button) (keyboard tgbotapi.ReplyKeyboardMarkup) {
	row := make([]tgbotapi.KeyboardButton, 0)
	rows := make([][]tgbotapi.KeyboardButton, 0)
	for i, b := range buttons {
		butt := tgbotapi.NewKeyboardButton(b.String())
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

func (list *UserList) PopUser(id int64) (User *user, f bool) {
	store := list.PopStore()
	User, ok := store[id]
	list.UnlockStore()

	return User, ok
}

func (list *UserList) PushUser(User *user) {
	store := list.PopStore()
	store[User.User_id] = User
	list.UnlockStore()
}

//Syncronise Userlist with DB info every 10 minutes
func updateUserList(botID int64) {
	upTime := time.Minute * 10
	for {
		store := Users.PopStore()
		newUsers, err := getUsersForBot(botID)
		if err != nil {
			log.Println(err)
			Users.UnlockStore()
			time.Sleep(upTime)
			continue
		}
		upUsers := make(map[int64]*user)
		for _, u := range newUsers {
			if user, ok := store[u.User_id]; ok {
				u.cmd = user.cmd
				u.prevCmd = user.prevCmd
				u.mutex = user.mutex

			} else {
				u.mutex = new(sync.Mutex)
			}
			tmp := u
			upUsers[u.User_id] = &tmp
		}
		Users.PushStore(upUsers)
		Users.UnlockStore()
		time.Sleep(upTime)
	}
}

func updateChats(bot *syncBot, botID int64) {

	for {
		chats, err := getChatsForBot(botID)
		if err == nil {
			for i, c := range chats {
				updateChatInfo(bot, c.Chat_id, &c)
				chats[i] = c
			}
		} else {
			log.Println(err)
		}
		time.Sleep(time.Minute * 10)
	}
}

func updateChatInfo(bot *syncBot, chatID int64, Chat *chat) {
	config := tgbotapi.ChatInfoConfig{}
	config.ChatID = chatID
	newChat, err := bot.GetChat(config)
	if err == nil {
		newDep := getDepartment(newChat.Title)
		if Chat.Title != newChat.Title || Chat.Department != newDep {
			Chat.Department = newDep
			Chat.Title = newChat.Title
			saveChat(*Chat)
		}
	} else {
		log.Println(err)
	}

}

func getChatsForBot(botID int64) (chats []chat, err error) {
	url := fmt.Sprintf("http://tg_cache:3334/chat/list/%d", botID)
	resp, err := http.Get(url) //TODO change to config parse
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

func getUsersForBot(botID int64) (Users []user, err error) {
	url := fmt.Sprintf("http://tg_cache:3334/user/list/%d", botID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		err = fmt.Errorf("status code = %d\n%s", response.StatusCode, string(responseBody))
		return
	}
	err = json.Unmarshal(responseBody, &Users)
	return
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

func newChat(c tgbotapi.Chat, bot_id int64) (Chat chat) {
	Chat.Bot_id = bot_id
	Chat.Chat_id = c.ID
	Chat.Department = getDepartment(c.Title)
	Chat.Title = c.Title
	switch c.Type {
	case "private":
		Chat.Type = 1
	default:
		Chat.Type = 2
	}
	return Chat
}

func (bot *syncBot) syncSend(value tgbotapi.Chattable) (response *tgbotapi.APIResponse, err error) {
	time.Sleep(time.Millisecond * 200)
	bot.mutex.Lock()
	response, err = bot.Request(value)
	switch response.ErrorCode {
	case 400:
		if response.Parameters != nil && response.Parameters.MigrateToChatID != 0 {
			x, err1 := convChattable(value)
			if err1 == nil {
				updateChatID(x.getBaseChat().ChatID, response.Parameters.MigrateToChatID, bot)
				bot.mutex.Unlock()
				x.changeChatID(response.Parameters.MigrateToChatID)
				response, err = bot.syncSend(x.getChattable())
				bot.mutex.Lock()
			}
		}
	case 403: //TODO finish this
		if strings.Contains(response.Description, "chat was deleted") ||
			strings.Contains(response.Description, "bot was kicked") {
			x, err1 := convChattable(value)
			if err1 == nil {
				b := x.getBaseChat()
				deleteChat(chat{Bot_id: bot.Self.ID, Chat_id: b.ChatID})
			}
		}
	}
	bot.mutex.Unlock()
	return
}

func newSyncBot() (bot *syncBot) {
	bot = new(syncBot)
	bot.mutex = new(sync.Mutex)
	return bot
}

func getDepartment(title string) (dep string) {
	dep = regexp.MustCompile(`ОС$`).FindString(title)

	if dep == "ОС"{
		return dep
	}
	return "ОЗ"
}

func (b Button) String() string {
	return string(b)
}

/*Checking report lifetime, if it is not beeing closed for 30 minutes, it's being closed automaticaly
checks lifetime every 3 minutes*/
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

func saveChat(Chat chat) (err error) {
	url := "http://tg_cache:3334/chat/add/"
	js, err := json.Marshal(Chat)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(js))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	resText, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(string(resText))
	}
	return nil
}

func convChattable(c tgbotapi.Chattable) (ret tg_chattable, err error) {
	if x, ok := c.(tgbotapi.MessageConfig); ok {
		ret.message = &x
		ret.chattable = x
		return
	} else if x, ok := c.(tgbotapi.AudioConfig); ok {
		ret.audio = &x
		ret.chattable = x
		return
	} else if x, ok := c.(tgbotapi.LocationConfig); ok {
		ret.location = &x
		ret.chattable = x
		return
	} else if x, ok := c.(tgbotapi.PhotoConfig); ok {
		ret.photo = &x
		ret.chattable = x
		return
	} else if x, ok := c.(tgbotapi.VideoConfig); ok {
		ret.video = &x
		ret.chattable = x
		return
	} else {
		return ret, errors.New("error converting")
	}
}

func (x *tg_chattable) getBaseChat() (Chat tgbotapi.BaseChat) {
	switch {
	case x.audio != nil:
		Chat = x.audio.BaseChat
	case x.location != nil:
		Chat = x.location.BaseChat
	case x.message != nil:
		Chat = x.message.BaseChat
	case x.photo != nil:
		Chat = x.photo.BaseChat
	case x.video != nil:
		Chat = x.video.BaseChat
	default:
		return Chat
	}
	return
}

func (x *tg_chattable) changeChatID(id int64) (Chat tgbotapi.BaseChat, err error) {
	switch {
	case x.audio != nil:
		x.audio.ChatID = id
		x.chattable = x.audio
	case x.location != nil:
		x.location.ChatID = id
		x.chattable = x.location
	case x.message != nil:
		x.message.ChatID = id
		x.chattable = x.message
	case x.photo != nil:
		x.photo.ChatID = id
		x.chattable = x.photo
	case x.video != nil:
		x.video.ChatID = id
		x.chattable = x.video
	default:
		return Chat, errors.New("error")
	}
	return
}

func (x *tg_chattable) getChattable() (ch tgbotapi.Chattable) {
	switch {
	case x.audio != nil:
		return *x.audio
	case x.location != nil:
		return *x.location
	case x.message != nil:
		return *x.message
	case x.photo != nil:
		return *x.photo
	case x.video != nil:
		return *x.video
	default:
		return nil
	}
}

func deleteChat(Chat chat) error {
	url := fmt.Sprintf("http://tg_cache:3334/chat/%d/%d", Chat.Bot_id, Chat.Chat_id)
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	resText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(string(resText))
	}
	return nil
}

func getKeyboard(u user, board Board) (keyboard []Button) {
	switch u.Rang {
	case AdminLvl:
		switch u.Department {
		case OS.String():
			switch board {
			case MainMenuBoard:
				return OSAdminMenu
			case DistribBoard:
				return OSDistribMenu
			case DepartmentSelectBoard:
				return OSSelectDepMenu
			}
		case OZ.String():
			switch board {
			case MainMenuBoard:
				return OZAdminMenu
			case DistribBoard:
				return OZDistribMenu
			case DepartmentSelectBoard:
				return OZSelectDepMenu
			}
		default:
			switch board {
			case MainMenuBoard:
				log.Println("supe")
				return ITAdminMenu
			case DistribBoard:
				return ITDistribMenu
			case DepartmentSelectBoard:
				return ITSelectDepMenu
			}
		}
	case SuperLvl:
		switch board {
		case MainMenuBoard:
			return SuperUserMenu
		case DistribBoard:
			return SuperDistribMenu
		case DepartmentSelectBoard:
			return SuperSelectDepMenu
		}
	case AnyLvl:
		switch u.Department {
		case OS.String():
			return OSUserMenu
		case OZ.String():
			return OZUserMenu
		}
	}
	return
}

/*Once pop is called, mutex getting locked,
don't forget to unlock
*/
func (list *UserList) PopStore() map[int64]*user {
	list.mutex.Lock()
	return list.store
}

func (list *UserList) PushStore(newSt map[int64]*user) {
	list.store = newSt
}

func (list *UserList) UnlockStore() {
	list.mutex.Unlock()
}

func (User *user) Block() {
	User.mutex.Lock()
}

func (User *user) Unblock() {
	User.mutex.Unlock()
}

func updateChatID(oldID, newID int64, bot *syncBot) {
	config := tgbotapi.ChatInfoConfig{}
	config.ChatID = newID
	NewTgChat, err := bot.GetChat(config)
	if err != nil {
		log.Println(err)
		return
	}
	config.ChatID = oldID
	OldTgChat, err := bot.GetChat(config)
	if err != nil {
		log.Println(err)
		return
	}
	err = deleteChat(newChat(OldTgChat, bot.Self.ID))
	if err != nil {
		log.Println(err)
		return
	}
	newC := newChat(NewTgChat, bot.Self.ID)
	saveChat(newC)
}

func saveUser(User user) (err error) {
	url := "http://tg_cache:3334/user/add"
	data, err := json.Marshal(User)
	if err != nil {
		return
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		res, _ := ioutil.ReadAll(response.Body)
		err = errors.New(string(res))
	}
	return
}

func (bot *syncBot) GetChat(config tgbotapi.ChatInfoConfig) (Chat tgbotapi.Chat, err error) {
	resp, err := bot.syncSend(config)
	if err != nil {
		return tgbotapi.Chat{}, err
	}

	err = json.Unmarshal(resp.Result, &Chat)

	return Chat, err
}
