package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type state int32

// const (
// 	mainMenu      = 0
// 	distrib       = 1
// 	distribGetMsg = 111
// 	hello         = 2
// )
var (
	keyboards = map[string][]string{
		"superUserMenu": {"Рассылка"},
		"adminMenu":     {"Рассылка"},
		"userMenu":      {"Привет"},
		"distribMenu":   {message, document, photo, video, "Главное меню", "Назад"},
	}
	superLvl uint8  = 0
	adminLvl uint8  = 1
	anyLvl   uint8  = 2
	message  string = "Сообщение"
	document string = "Документ"
	photo    string = "Фото"
	video    string = "Видео"
	distrib  string = "Рассылка"
	mainMenu string = "Главное меню"
)

type mainMenuCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       string
	executable bool
}

type distribCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       string
	executable bool
}

type backCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       string
	executable bool
}

type unknownCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       string
	executable bool
}

type Cmd interface {
	init(*user, string)
	exec(*tgbotapi.BotAPI, *tgbotapi.Update) error
	String() string
	setName(string)
	copy() Cmd
}

type chat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	BotID int64  `json:"botID"`
	Type  uint8  `json:"type"`
}

type user struct {
	ID         int64  `json:"id"`
	BotID      int64  `json:"botId"`
	Username   string `json:"username"`
	Firstname  string `json:"firstname,omitempty"`
	Rang       uint8  `json:"rang"` //0 - superUser, 1 - admin
	Department string `json:"department"`
	prevCmd    Cmd
	cmd        Cmd
}

type callbackJSON struct {
	Type string `json:"info"`
	Info string `json:"text"`
}
