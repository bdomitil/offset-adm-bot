package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *addUserCmd) init(u *user, cmd Button) {
	c.name = cmd
	c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*u, MainMenuBoard))
	c.executable = true
	if u.Rang > SuperLvl {
		c.executable = false
	}
	c.level = SuperLvl
	c.state = Processing
}

func (c *addUserCmd) exec(bot *syncBot, u tgbotapi.Update) (err error) {
	var msg tgbotapi.MessageConfig
	msg.ChatID = u.Message.From.ID
	User, _ := Users.PopUser(u.SentFrom().ID)
	switch c.name {
	case AddUser:
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(getKeyboard(*User, DepartmentSelectBoard))
		msg.Text = "Выберите отдел пользователя"
	case OS:
		msg.Text = "Пожалуйста отправьте профиль пользователя из " + OS.String()
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(getKeyboard(*User, MainMenuBoard))
		c.state = W8message
	case OZ:
		msg.Text = "Пожалуйста отправьте профиль пользователя из " + OZ.String()
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(getKeyboard(*User, MainMenuBoard))
		c.state = W8message
	case Stop:
		msg.ReplyMarkup = NewResizeOneTimeReplyKeyboard(getKeyboard(*User, MainMenuBoard))
		c.state = Closed
		err = addUser(u, bot, User.prevCmd.String())
		if err != nil {
			msg.Text = err.Error()
		} else {
			msg.Text = fmt.Sprintf("Пользователь %s успешно добавлен", u.Message.Contact.FirstName)
		}
	}
	_, err = bot.syncSend(msg)
	return
}

func (c *addUserCmd) String() string {
	return string(c.name)
}

func (c *addUserCmd) setName(n string) {
	c.name = Button(n)
}

func (c *addUserCmd) copy() (copy Cmd) {
	cp := new(addUserCmd)
	cp.level = c.level
	cp.executable = c.executable
	cp.keyboard = c.keyboard
	cp.state = c.state
	cp.name = c.name
	return cp
}

func (c *addUserCmd) getState() State {
	return c.state
}

func (c *addUserCmd) setState(s State) {
	c.state = s
}

func (c *unknownCmd) init(u *user, cmd Button) {
	c.name = "Error"
	c.level = AnyLvl
	c.executable = true
	c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*u, MainMenuBoard))
	c.state = Processing
}

func (c *mainMenuCmd) init(u *user, cmd Button) {
	c.name = MainMenu
	c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*u, MainMenuBoard))
	c.executable = true
	c.level = AnyLvl
	c.state = Processing
}

func (c *distribCmd) init(u *user, cmd Button) {
	c.name = cmd
	c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*u, MainMenuBoard))
	if u.Rang <= 1 {
		c.executable = true
	}
	c.level = AdminLvl
	c.state = Processing
}

func (c *backCmd) init(u *user, cmd Button) {
	c.name = "Назад"
	c.level = AnyLvl
	c.state = Processing
}

func (c *unknownCmd) exec(bot *syncBot, u tgbotapi.Update) (err error) {
	var msg tgbotapi.MessageConfig
	msg.ChatID = u.SentFrom().ID
	msg.Text = "Неизвестная команда"
	msg.ReplyMarkup = c.keyboard
	_, err = bot.syncSend(msg)
	c.state = Closed
	return
}

func (c *mainMenuCmd) exec(bot *syncBot, u tgbotapi.Update) (err error) {
	var msg tgbotapi.MessageConfig
	msg.ChatID = u.SentFrom().ID
	msg.Text = MainMenu.String()
	msg.ReplyMarkup = c.keyboard
	_, err = bot.syncSend(msg)
	c.state = Processing
	return
}
func (c *backCmd) exec(bot *syncBot, u tgbotapi.Update) (err error) {

	var msg tgbotapi.MessageConfig
	msg.ChatID = u.SentFrom().ID
	c.state = Processing
	msg.Text = "Назад"
	// if

	return
}

func (c *distribCmd) setName(newName string) {
	c.name = Button(newName)
}
func (c *mainMenuCmd) setName(newName string) {
	c.name = Button(newName)
}
func (c *backCmd) setName(newName string) {
	c.name = Button(newName)
}
func (c *unknownCmd) setName(newName string) {
	c.name = Button(newName)
}

func (c *mainMenuCmd) copy() (copy Cmd) {
	x := new(mainMenuCmd)
	x.executable = c.executable
	x.keyboard = c.keyboard
	x.level = c.level
	x.name = c.name
	x.state = c.state
	return x
}

func (c *distribCmd) copy() (copy Cmd) {
	x := new(distribCmd)
	x.executable = c.executable
	x.keyboard = c.keyboard
	x.level = c.level
	x.name = c.name
	x.state = c.state
	return x
}

func (c *backCmd) copy() (copy Cmd) {
	x := new(backCmd)
	x.executable = c.executable
	x.keyboard = c.keyboard
	x.level = c.level
	x.name = c.name
	x.state = c.state
	return x
}

func (c *unknownCmd) copy() (copy Cmd) {
	x := new(unknownCmd)
	x.executable = c.executable
	x.keyboard = c.keyboard
	x.level = c.level
	x.name = c.name
	x.state = c.state
	return x
}

func (c *mainMenuCmd) getState() State {
	return c.state
}

func (c *mainMenuCmd) setState(s State) {
	c.state = s
}

func (c *unknownCmd) getState() State {
	return c.state
}

func (c *unknownCmd) setState(s State) {
	c.state = s
}

func (c *backCmd) getState() State {
	return c.state
}

func (c *backCmd) setState(s State) {
	c.state = s
}

func (c *distribCmd) getState() State {
	return c.state
}

func (c *distribCmd) setState(s State) {
	c.state = s
}

func resend_as_distrib(bot *syncBot, u tgbotapi.Update) (err error) {
	chats, err := getChatsForBot(bot.Self.ID)
	user, _ := Users.PopUser(u.SentFrom().ID)
	if err != nil {
		return
	}
	for _, c := range chats {
		if user.Department != c.Department && user.Department != "ИТ" {
			continue
		}
		var chattable tgbotapi.Chattable
		switch {
		case u.Message.Document != nil:
			chattable = tgbotapi.NewDocument(c.Chat_id, tgbotapi.FileID(u.Message.Document.FileID))
		case u.Message.Photo != nil:
			chattable = tgbotapi.NewPhoto(c.Chat_id, tgbotapi.FileID(u.Message.Photo[len(u.Message.Photo)-1].FileID))
		case u.Message.Video != nil:
			chattable = tgbotapi.NewVideo(c.Chat_id, tgbotapi.FileID(u.Message.Video.FileID))
		case len(u.Message.Text) > 0:
			chattable = tgbotapi.NewMessage(c.Chat_id, u.Message.Text)
		}
		_, err = bot.syncSend(chattable)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func (c *distribCmd) exec(bot *syncBot, u tgbotapi.Update) (err error) {
	var msg tgbotapi.MessageConfig
	msg.ChatID = u.SentFrom().ID
	user, _ := Users.PopUser(u.SentFrom().ID)
	switch c.name {
	case Distrib:
		msg.Text = "Пожалуйста выберите тип сообщение для массовой рассылки"
		c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*user, DistribBoard))
		c.state = Processing
	case Photo:
		msg.Text = "Пожалуйста отправьте мне фото без ТЕКСТА К НЕМУ!"
		c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*user, MainMenuBoard))
		c.state = W8message
	case Document:
		msg.Text = "Пожалуйста отправьте мне докумен без ТЕКСТА К НЕМУ!"
		c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*user, MainMenuBoard))
		c.state = W8message
	case Video:
		msg.Text = "Пожалуйста отправьте мне видео без ТЕКСТА К НЕМУ!"
		c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*user, MainMenuBoard))
		c.state = W8message
	case Message:
		msg.Text = "Пожалуйста отправьте мне сообщение без МЕДИА файлов!"
		c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*user, MainMenuBoard))
		c.state = W8message
	case Stop:
		err = resend_as_distrib(bot, u)
		if err != nil {
			log.Println(err)
		}
		msg.Text = "Рассылка окончена"
		c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*user, DistribBoard))
		c.state = Closed
	default:
		log.Println("unknown distr")
		c.keyboard = NewResizeOneTimeReplyKeyboard(getKeyboard(*user, MainMenuBoard))
		c.state = Closed
	}
	msg.ReplyMarkup = c.keyboard
	_, err = bot.syncSend(msg)
	return
}

func (c *unknownCmd) button() Button {
	return c.name
}
func (c *addUserCmd) button() Button {
	return c.name
}
func (c *mainMenuCmd) button() Button {
	return c.name
}
func (c *backCmd) button() Button {
	return c.name
}
func (c *distribCmd) button() Button {
	return c.name
}

func (c *unknownCmd) String() string {
	return string(c.name)
}

func (c *distribCmd) String() string {
	return string(c.name)
}

func (c *mainMenuCmd) String() string {
	return string(c.name)
}
func (c *backCmd) String() string {
	return string(c.name)
}

func (u *user) newCmd(cmd Button) (newCmd Cmd, err error) {

	switch {
	case isNil(u.prevCmd) || u.prevCmd.getState() != W8message:
		{
			switch cmd {
			case "/start":
				fallthrough
			case MainMenu:
				newCmd = new(mainMenuCmd)
			case Back:
				newCmd = new(backCmd)
			case Message:
				fallthrough
			case Document:
				fallthrough
			case Photo:
				fallthrough
			case Video:
				if !isNil(u.prevCmd) && u.prevCmd.button() != Distrib {
					newCmd = new(unknownCmd)
					break
				}
				fallthrough
			case Distrib:
				newCmd = new(distribCmd)
			case OS:
				fallthrough
			case OZ:
				fallthrough
			case AddUser:
				newCmd = new(addUserCmd)
			default:
				newCmd = new(unknownCmd)
			}
		}
	case u.prevCmd.getState() == W8message:
		{

			if _, ok := u.prevCmd.(*distribCmd); ok {
				cmd = Stop
				newCmd = new(distribCmd)
			} else if _, ok := u.prevCmd.(*addUserCmd); ok {
				cmd = Stop
				newCmd = new(addUserCmd)
			} else {
				newCmd = new(unknownCmd)
			}
		}
	}
	newCmd.init(u, cmd)
	return newCmd, nil
}

func addUser(update tgbotapi.Update, bot *syncBot, dep string) (err error) {

	member, err := getChatMember(update.FromChat().ID, update.Message.Contact.UserID, bot)
	if err != nil {
		return
	}
	User := user{
		User_id:    member.User.ID,
		Bot_id:     bot.Self.ID,
		Firstname:  member.User.FirstName,
		Rang:       AdminLvl, //TODO Hardcoded
		Department: dep,
		Username:   member.User.UserName,
	}
	return saveUser(User)
}

func getChatMember(ChatID, UserID int64, bot *syncBot) (member tgbotapi.ChatMember, err error) {
	cf := tgbotapi.GetChatMemberConfig{}
	cf.UserID = UserID
	cf.ChatID = ChatID
	member, err = bot.GetChatMember(cf)
	if err != nil {
		return
	}
	return
}
