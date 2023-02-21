package sender

import (
	"context"
	"github.com/COTBU/notifier/config"
	"github.com/COTBU/notifier/pkg/model"
	"github.com/COTBU/notifier/service/sender/email"
	"github.com/COTBU/notifier/service/sender/telegram"
)

type Sender struct {
	config *config.Config
}

func New(config *config.Config) *Sender {
	return &Sender{config: config}
}

func (s *Sender) ProcessMessage(notification model.Notification) error {
	switch notification.Type {
	case model.TelegramNotification:
		if err := s.SendTelegramMessage(notification); err != nil {
			return err
		}
	default:
		if err := s.SendEmail(notification); err != nil {
			return err
		}
	}

	return nil
}

func (s *Sender) SendEmail(notification model.Notification) error {
	client, err := email.NewClient(s.config)
	if err != nil {
		return err
	}

	return client.
		SetSubject(notification.Subject).
		SetDestination(notification.Recipients).
		SendRich(string(notification.Body))
}

func (s *Sender) SendTelegramMessage(notification model.Notification) error {
	client, err := telegram.NewClient(s.config)
	if err != nil {
		return err
	}

	return client.Send(
		context.Background(),
		notification.Subject,
		string(notification.Body),
	)
}
