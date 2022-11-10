package main

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type State int32
type Button string
type Board int32

const ()

var (
	MainMenuBoard         Board    = 0
	DistribBoard          Board    = 1
	DepartmentSelectBoard Board    = 2
	SuperUserMenu         []Button = []Button{Distrib, AddUser}
	OSAdminMenu           []Button = []Button{Distrib}
	OZAdminMenu           []Button = []Button{Distrib}
	ITAdminMenu           []Button = []Button{Distrib}
	OZUserMenu            []Button = []Button{Hi}
	OSUserMenu            []Button = []Button{Hi}
	SuperDistribMenu      []Button = []Button{Message, Document, Photo, Video, MainMenu, Back}
	OSDistribMenu         []Button = []Button{Message, Document, Photo, Video, MainMenu, Back}
	ITDistribMenu         []Button = []Button{Message, Document, Photo, Video, MainMenu, Back}
	OZDistribMenu         []Button = []Button{Message, Document, Photo, Video, MainMenu, Back}
	ODistribMenu          []Button = []Button{Message, Document, Photo, Video, MainMenu, Back}
	ITSelectDepMenu       []Button = []Button{OS, OZ, MainMenu}
	OSSelectDepMenu       []Button = []Button{OS, OZ, MainMenu}
	OZSelectDepMenu       []Button = []Button{OS, OZ, MainMenu}
	SuperSelectDepMenu    []Button = []Button{OS, OZ, MainMenu}
	SuperLvl              uint8    = 0
	AdminLvl              uint8    = 1
	AnyLvl                uint8    = 2
	Message               Button   = "Сообщение"
	Hi                    Button   = "Привет"
	AddUser               Button   = "Add user"
	Document              Button   = "Документ"
	Photo                 Button   = "Фото"
	Video                 Button   = "Видео"
	Distrib               Button   = "Рассылка"
	MainMenu              Button   = "Главное меню"
	Back                  Button   = "Назад"
	Stop                  Button   = "Stop"
	OS                    Button   = "ОС"
	OZ                    Button   = "ОЗ"
	W8message             State    = 0
	Processing            State    = 1
	Closed                State    = 2
)

type UserConfig struct {
	ID         int64
	Firstname  string
	Department string
	Lastname   string
	Rang       int32
}

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

type syncBot struct {
	*tgbotapi.BotAPI
	mutex sync.Mutex
}

type mainMenuCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       Button
	executable bool
	state      State
}

type addUserCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       Button
	executable bool
	state      State
}

type distribCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       Button
	executable bool
	state      State
}

type backCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       Button
	executable bool
	state      State
}

type unknownCmd struct {
	keyboard   tgbotapi.ReplyKeyboardMarkup
	level      uint8 //0 - the highest level need to be executed, 2 - anyone can execute
	name       Button
	executable bool
	state      State
}

type Cmd interface {
	init(*user, Button)
	exec(*syncBot, *tgbotapi.Update) error
	String() Button
	setName(string)
	copy() Cmd
	getState() State
	setState(s State)
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
