package main

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	superLvl   uint8  = 0
	adminLvl   uint8  = 1
	anyLvl     uint8  = 2
	message    string = "Сообщение"
	document   string = "Документ"
	photo      string = "Фото"
	video      string = "Видео"
	distrib    string = "Рассылка"
	mainMenu   string = "Главное меню"
	w8message  state  = 0
	processing state  = 1
	closed     state  = 2
)

type reportList struct {
	store map[int64]report
	mutex sync.Mutex
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
	chat_id       int64
	channel_name  string
	creation_time time.Time
	isFilled      bool
	openMsgID     int
}

type tg_chattable struct {
	message   *tgbotapi.MessageConfig
	photo     *tgbotapi.PhotoConfig
	video     *tgbotapi.VideoConfig
	location  *tgbotapi.LocationConfig
	audio     *tgbotapi.AudioConfig
	chattable tgbotapi.Chattable
}

type syncBot struct {
	*tgbotapi.BotAPI
	mutex *sync.Mutex
}

type mainMenuCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       string
	executable bool
	state      state
}

type distribCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       string
	executable bool
	state      state
}

type backCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       string
	executable bool
	state      state
}

type unknownCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       string
	executable bool
	state      state
}

type Cmd interface {
	init(*user, string)
	exec(*syncBot, *tgbotapi.Update) error
	String() string
	setName(string)
	copy() Cmd
	getState() state
	setState(s state)
}

type chat struct {
	Chat_id    int64  `json:"chat_id" gorm:"column:chat_id"`
	Title      string `json:"title" gorm:"column:title"`
	Bot_id     int64  `json:"bot_id" gorm:"column:bot_id"`
	Type       uint8  `json:"type" gorm:"column:type"`
	Department string `json:"department" gorm:"column:department"`
}

type user struct {
	User_id    int64  `json:"user_id" gorm:"column:user_id"`
	Bot_id     int64  `json:"bot_id" gorm:"column:bot_id"`
	Username   string `json:"username" gorm:"column:username"`
	Firstname  string `json:"firstname,omitempty" gorm:"column:firstname"`
	Rang       uint8  `json:"rang" gorm:"column:rang"`
	Department string `json:"department" gorm:"column:department"`
	prevCmd    Cmd
	cmd        Cmd
}

type callbackJSON struct {
	Type string `json:"info"`
	Info string `json:"text"`
}
