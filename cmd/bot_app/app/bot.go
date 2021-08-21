package app

import (
	"bytes"
	"convert-song-tg-bot/pkg/config"
	"fmt"
	"github.com/NicoNex/echotron/v3"
	"net/url"
	"strings"
	"text/template"
)

type bot struct {
	chatID int64
	echotron.API
}

func (b *bot) SendSimpleMess(mess string) error {
	_, err := b.SendMessage(mess, b.chatID, nil)
	return err
}

func (b *bot) SendHtmlMess(format *template.Template, fields EscapeCharsMarkdown) error {
	opts := &echotron.MessageOptions{
		ParseMode: echotron.HTML,
	}
	fields.CheckFields(EscapeChars)
	buf := &bytes.Buffer{}
	fmt.Printf("Fields %#v\n tmpl %#v\n writer %#v", fields, format, buf)
	err := format.Execute(buf, fields)
	if err != nil {
		return err
	}
	b.SendMessage(fmt.Sprint(buf), b.chatID, opts)
	return nil
}

func EscapeChars(line string) string {
	line = strings.ReplaceAll(line, "<", "&lt;")
	line = strings.ReplaceAll(line, ">", "&gt")
	return strings.ReplaceAll(line, "&", "&amp")
}

func (b *bot) SendStickerMess(id string) error {
	_, err := b.SendSticker(id, b.chatID, nil)
	return err
}

type UserStruct struct {
	*echotron.User
}

func (u UserStruct) getFirstName() string {
	return u.FirstName
}

func NewBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(config.GetEnv("TELEGRAM_API_TOKEN", "")),
	}
}

func inputRouter(b *bot, update *echotron.Update, arHandlersCheck []isHandler) bool {
	HandlerFounded := false
	for _, item := range arHandlersCheck {
		if item(b, update) {
			HandlerFounded = true
			break
		}

	}
	return HandlerFounded
}

//Step for router input messages
//Returns true, if check passed
type isHandler func(b *bot, update *echotron.Update) bool

func isCommand(b *bot, update *echotron.Update) bool {
	if update.Message.Text != "" {

		strAr := strings.SplitN(update.Message.Text, "", 2)
		if strAr[0] == "/" && len(strAr) > 1 {
			fmt.Println("command detected")
			var args string
			commandAr := strings.SplitN(strAr[1], " ", 2)

			if len(commandAr) > 1 {
				args = commandAr[1]
			}
			User := &UserStruct{User: update.Message.From}
			hanlder := &CommandHandler{
				BaseHandler: BaseHandler{B: b, Input: update.Message.Text, U: User},
				Command:     commandAr[0],
				Argument:    args,
			}
			err := hanlder.Execute()
			if err != nil {
				//Todo need log
				fmt.Println("Error in command")
				b.SendMessage("Shit happend, try else", b.chatID, nil)

			}
			return true
		}
	}
	fmt.Println("not command")
	return false
}

func isSticker(b *bot, update *echotron.Update) bool {
	//Short handler
	if update.Message.Sticker != nil {
		b.SendStickerMess("CAACAgIAAxkBAAIBcGEXamEj-fQAAbX5NBwRebgO3sVb7QACGwIAAtzyqwcxrg0ZUPeeBiAE")
		return true
	}
	return false
}

func isLink(b *bot, update *echotron.Update) bool {
	fmt.Println("is Link started")
	if update.Message.Text != "" {
		link, err := url.Parse(update.Message.Text)

		if err != nil {
			return false
		}
		if link.Scheme != "" && link.Host != "" {
			fmt.Println("handler link go!")
			User := &UserStruct{User: update.Message.From}
			hanlder := &LinkHandler{
				BaseHandler: BaseHandler{B: b, Input: update.Message.Text, U: User},
				InputLink:   update.Message.Text,
			}
			hanlder.Execute()
			return true
		}
	}
	return false
}

func (b *bot) Update(update *echotron.Update) {

	arHandlersCheck := []isHandler{
		isHandler(isCommand),
		isHandler(isLink),
		isHandler(isSticker),
	}
	if !inputRouter(b, update, arHandlersCheck) {
		b.SendMessage("Я не понимаю. Попробуй /help", b.chatID, nil)
	}
}
