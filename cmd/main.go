package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/wizardloong/botposter/pkg/handler"
	"github.com/wizardloong/botposter/pkg/service"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	telegramToken := os.Getenv("BOT_TOKEN")

	bot, err := tb.NewBot(tb.Settings{
		Token:  telegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("error new bot: %s", err.Error())
		return
	}

	services := service.NewServices()
	handler := handler.NewHandler(services, bot)

	handler.RunBot()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
