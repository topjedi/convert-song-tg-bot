package app

import (
	"fmt"
	"github.com/NicoNex/echotron/v3"
	"net/url"
	"strings"
	"text/template"
)

//----Parent of all handlers
type BaseHandler struct {
	U     UserI
	B     Sender
	Input string
}

type UserI interface {
	getFirstName() string
}

//=====COMMAND=====
type CommandHandler struct {
	BaseHandler
	Command  string
	Argument string
}

func (command *CommandHandler) Execute() error {
	switch command.Command {
	case "start":
		mess := "Привет, " + command.U.getFirstName() + "! Я помогу найти нужную песню на разных стримингах, просто пришли ссылку"
		command.B.SendSimpleMess(mess)
		//b.SendMessage("Привет, "+update.Message.User.FirstName+"! Я помогу найти нужную песню на разных стримингах, просто пришли ссылку", b.chatID, nil)
		command.B.SendStickerMess("CAACAgIAAxkBAAPBYQGg456U9Lgyk071EF7zAAH0HA6KAAIQAgAC3PKrB5L_imUjZzVnIAQ")
	case "help":
		mess := "/web <link_url> - Результаты в виде веб страницы\n/supported - поддерживаемые платформы\n<link_url> - все что удастся найти"
		command.B.SendSimpleMess(mess)
	case "supported":
		mess := "- Amazon Music\n- Apple Music\n- Audius\n- Deezer\n- iTunes\n- Napster\n- Pandora\n- SoundCloud\n- Spinrilla\n- Spotify\n- TIDAL\n- Yandex.Music\n- YouTube (videos)\n- YouTube Music"
		command.B.SendSimpleMess(mess)
	case "web":

		link, err := url.Parse(command.Argument)
		if err != nil {
			command.B.SendSimpleMess("Некорректная ссылка")
			return nil
		}
		fmt.Printf("ARGUMENT: %#v\n", command.Argument)
		fmt.Printf("WEB LINK: %#v\n", link)
		if link.Scheme != "" && link.Host != "" {
			song, err := getSong(command.Argument)
			fmt.Printf("SONG %#v\n", song)
			if err != nil {
				command.B.SendSimpleMess("Не удалось получить результат")
				return nil
			} else {
				tpl := `<a href="{{.Web}}">🎤<b>{{.Title}}</b> - {{.Artist}}</a>`
				t, e := template.New("SongResponceWeb").Parse(tpl)
				if e != nil {
					fmt.Printf("Error to parse templ: %#v", e)
					return e
				}
				return command.B.SendHtmlMess(t, song)
			}
		} else {
			command.B.SendSimpleMess("Некорректная ссылка")
			return nil
		}
	default:
		mess := "Я не знаю такой команды. Если нужна помощь отправь /help"
		command.B.SendSimpleMess(mess)
		command.B.SendStickerMess("CAACAgIAAxkBAAIBbmEXaZ6gipTVFeQYuNmfm96iTghVAAIkAgAC3PKrB8XVblz6HdRqIAQ")
	}
	return nil
}

//=====END COMMAND=====

//=====LINK=====
type LinkHandler struct {
	BaseHandler
	InputLink string
}

func (handler *LinkHandler) Execute() error {
	song, err := getSong(handler.InputLink)
	if err != nil {
		return err
	}
	tpl := `🎤<b>{{.Title}}</b> - {{.Artist}}<a href="{{.Pic}}">.</a>
Вот что мне удалось найти:
{{range .Links}}
<a href="{{.Url}}">&lt;{{.Name}}&gt</a>
{{end}}`
	t, e := template.New("SongResponce").Parse(tpl)
	if e != nil {
		fmt.Printf("Error to parse templ: %#v", e)
	}
	for i := range song.Links {
		song.Links[i].Name = strings.Title(song.Links[i].Name)
	}
	return handler.B.SendHtmlMess(t, song)
}

type Song struct {
	Web    string
	Title  string
	Artist string
	Type   string
	Pic    string
	Links  []LinkSong
}

type LinkSong struct {
	Name string
	Url  string
}

func (song *Song) CheckFields(escape EscapeCharsField) {
	song.Title = escape(song.Title)
	song.Artist = escape(song.Artist)
	song.Pic = escape(song.Pic)
	for _, val := range song.Links {
		val.Name = escape(val.Name)
		val.Url = escape(val.Url)
	}
}

//=====END LINK=====

type StickerHandler struct {
	BaseHandler
	InputSticker echotron.Sticker
}
type Sender interface {
	SendSimpleMess(mess string) error
	SendHtmlMess(*template.Template, EscapeCharsMarkdown) error
	SendStickerMess(id string) error
}

type EscapeCharsField func(string) string

type EscapeCharsMarkdown interface {
	// CheckFields contains exampleStruct.fieldN = EscapeCharsField(exampleStruct.fieldN)
	CheckFields(EscapeCharsField)
}
