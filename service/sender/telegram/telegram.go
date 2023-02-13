package telegram

import (
	"SOTBI/notifier/config"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
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
