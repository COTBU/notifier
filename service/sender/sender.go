package sender

import (
	"github.com/COTBU/notifier/config"
	"github.com/COTBU/notifier/pkg/model"
	"github.com/COTBU/notifier/service/sender/email"
	"github.com/COTBU/notifier/service/sender/telegram"
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
