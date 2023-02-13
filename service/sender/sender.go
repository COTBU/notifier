package sender

import (
	"SOTBI/notifier/config"
	"SOTBI/notifier/model"
	"SOTBI/notifier/service/sender/email"
	"SOTBI/notifier/service/sender/telegram"
)

type Sender struct {
	config *config.Config
}

func (s *Sender) ProcessMessage(notification model.Notification) error {
	switch notification.Type {
	case model.TelegramNotification:
	default:

	}

	return nil
}

func (s *Sender) SendEmail() error {
	client, err := email.NewClient(s.config)
	if err != nil {
		return err
	}
	_ = client
	return nil
}

func (s *Sender) SendTelegramMessage() error {
	client, err := telegram.NewClient(s.config)
	if err != nil {
		return err
	}
	_ = client
	return nil
}
