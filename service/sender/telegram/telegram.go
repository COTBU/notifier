package telegram

import (
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"

	"github.com/COTBU/notifier/config"
)

func NewClient(appConfig *config.Config) (*telegram.Telegram, error) {
	telegramClient, err := telegram.New(appConfig.Telegram.Token)
	if err != nil {
		return nil, err
	}

	telegramClient.AddReceivers(appConfig.Telegram.Channel)

	notifyClient := notify.New()
	notifyClient.UseServices(telegramClient)

	return telegramClient, nil
}
