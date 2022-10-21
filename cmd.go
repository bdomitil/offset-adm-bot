package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *unknownCmd) init(u *user, cmd string) {
	c.name = "Error"
	c.level = anyLvl
	c.executable = true
	switch u.Rang {
	case adminLvl:
		c.keyboard = NewResizeOneTimeReplyKeyboard(keyboards["adminMenu"]...)
	case superLvl:
		c.keyboard = NewResizeOneTimeReplyKeyboard(keyboards["superUserMenu"]...)
	}
}

func (c *mainMenuCmd) init(u *user, cmd string) {
	c.name = mainMenu
	switch u.Rang {
	case adminLvl:
		c.keyboard = NewResizeOneTimeReplyKeyboard(keyboards["adminMenu"]...)
	case superLvl:
		c.keyboard = NewResizeOneTimeReplyKeyboard(keyboards["superUserMenu"]...)
	}
	c.executable = true
	c.level = anyLvl
}

func (c *distribCmd) init(u *user, cmd string) {
	c.name = cmd
	switch u.Rang {
	case adminLvl:
		fallthrough
	case superLvl:
		c.keyboard = NewResizeOneTimeReplyKeyboard(keyboards["distribMenu"]...)
	}
	if u.Rang <= 1 {
		c.executable = true
	}
	c.level = adminLvl
}

func (c *backCmd) init(u *user, cmd string) {
	c.name = "Назад"
	c.level = anyLvl
}

func (c *unknownCmd) exec(bot *tgbotapi.BotAPI, u *tgbotapi.Update) (err error) {
	var msg tgbotapi.MessageConfig
	msg.ChatID = u.FromChat().ID
	msg.Text = "Неизвестная команда"
	msg.ReplyMarkup = c.keyboard
	_, err = bot.Send(msg)
	return
}

func (c *mainMenuCmd) exec(bot *tgbotapi.BotAPI, u *tgbotapi.Update) (err error) {
	var msg tgbotapi.MessageConfig
	msg.ChatID = u.FromChat().ID
	msg.Text = mainMenu
	msg.ReplyMarkup = c.keyboard
	_, err = bot.Send(msg)
	return
}
func (c *backCmd) exec(bot *tgbotapi.BotAPI, u *tgbotapi.Update) (err error) {

	var msg tgbotapi.MessageConfig
	msg.ChatID = u.FromChat().ID
	msg.Text = "Назад"
	// if

	return
}

func (c *distribCmd) setName(newName string) {
	c.name = newName
}
func (c *mainMenuCmd) setName(newName string) {
	c.name = newName
}
func (c *backCmd) setName(newName string) {
	c.name = newName
}
func (c *unknownCmd) setName(newName string) {
	c.name = newName
}

func (c *mainMenuCmd) copy() (copy Cmd) {
	x := new(mainMenuCmd)
	x.executable = c.executable
	x.keyboard = c.keyboard
	x.level = c.level
	x.name = c.name
	return x
}

func (c *distribCmd) copy() (copy Cmd) {
	x := new(distribCmd)
	x.executable = c.executable
	x.keyboard = c.keyboard
	x.level = c.level
	x.name = c.name
	return x
}

func (c *backCmd) copy() (copy Cmd) {
	x := new(backCmd)
	x.executable = c.executable
	x.keyboard = c.keyboard
	x.level = c.level
	x.name = c.name
	return x
}

func (c *unknownCmd) copy() (copy Cmd) {
	x := new(unknownCmd)
	x.executable = c.executable
	x.keyboard = c.keyboard
	x.level = c.level
	x.name = c.name
	return x
}

func resend_as_distrib(bot *tgbotapi.BotAPI, u *tgbotapi.Update) (err error) {
	chats, err := getChatsForBot(bot.Self.ID)
	if err != nil {
		return
	}
	for _, c := range chats {
		var chattable tgbotapi.Chattable
		switch {
		case u.Message.Document != nil:
			chattable = tgbotapi.NewDocument(c.ID, tgbotapi.FileID(u.Message.Document.FileID))
		case u.Message.Photo != nil:
			chattable = tgbotapi.NewPhoto(c.ID, tgbotapi.FileID(u.Message.Photo[len(u.Message.Photo)-1].FileID))
		case u.Message.Video != nil:
			chattable = tgbotapi.NewVideo(c.ID, tgbotapi.FileID(u.Message.Video.FileID))
		case len(u.Message.Text) > 0:
			chattable = tgbotapi.NewMessage(c.ID, u.Message.Text)
		}
		_, err = bot.Send(chattable)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func (c *distribCmd) exec(bot *tgbotapi.BotAPI, u *tgbotapi.Update) (err error) {
	var msg tgbotapi.MessageConfig
	msg.ChatID = u.FromChat().ID
	switch c.name {
	case distrib:
		msg.Text = "Пожалуйста выберите тип сообщение для массовой рассылки"
		msg.ReplyMarkup = c.keyboard
	case photo:
		msg.Text = "Пожалуйста отправьте мне фото без ТЕКСТА К НЕМУ!"
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(mainMenu)
	case document:
		msg.Text = "Пожалуйста отправьте мне докумен без ТЕКСТА К НЕМУ!"
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(mainMenu)
	case video:
		msg.Text = "Пожалуйста отправьте мне видео без ТЕКСТА К НЕМУ!"
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(mainMenu)
	case message:
		msg.Text = "Пожалуйста отправьте мне сообщение без МЕДИА файлов!"
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(mainMenu)
	case "stop distrib":
		err = resend_as_distrib(bot, u)
		if err != nil {
			log.Fatalln(err)
		}
		msg.Text = "Рассылка оконченна"
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(mainMenu)
	default:
		log.Println("unknown distr")
	}
	_, err = bot.Send(msg)
	return
}

func (c *unknownCmd) String() string {
	return c.name
}

func (c *distribCmd) String() string {
	return c.name
}

func (c *mainMenuCmd) String() string {
	return c.name
}
func (c *backCmd) String() string {
	return c.name
}

func (u *user) newCmd(cmd string) (newCmd Cmd, err error) {

	switch cmd {
	case "/start":
		fallthrough
	case mainMenu:
		newCmd = new(mainMenuCmd)
	case "Назад":
		newCmd = new(backCmd)
	case message:
		fallthrough
	case document:
		fallthrough
	case photo:
		fallthrough
	case video:
		if !isNil(u.prevCmd) && u.prevCmd.String() != distrib {
			newCmd = new(unknownCmd)
			break
		}
		fallthrough
	case distrib:
		newCmd = new(distribCmd)
	default:
		if _, ok := u.prevCmd.(*distribCmd); ok { //TODO Make this part normal
			newCmd = new(distribCmd)
			cmd = "stop distrib"
		} else {
			newCmd = new(unknownCmd)
		}
	}
	newCmd.init(u, cmd)
	return newCmd, nil
}
