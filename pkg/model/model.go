package model

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type Notification struct {
	Type       uint
	Subject    string
	Body       []byte
	Recipients []string
}

const (
	TelegramNotification = iota + 1
)

func (n *Notification) Encode() ([]byte, error) {
	if n == nil {
		return nil, errors.New("encoding nil notification aborted")
	}
	buffer := &bytes.Buffer{}

	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(n); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (n *Notification) Decode(data []byte) error {
	if n == nil {
		return errors.New("decoding into nil notification aborted")
	}
	buffer := bytes.NewBuffer(data)

	decoder := gob.NewDecoder(buffer)

	return decoder.Decode(n)
}
