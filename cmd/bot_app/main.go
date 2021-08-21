package main

import (
	"convert-song-tg-bot/cmd/bot_app/app"
	"convert-song-tg-bot/pkg/config"
	"fmt"
	"github.com/NicoNex/echotron/v3"
)

func main() {
	dsp := echotron.NewDispatcher(config.GetEnv("TELEGRAM_API_TOKEN", ""), app.NewBot)
	fmt.Println(dsp.Poll())

}
