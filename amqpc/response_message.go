package amqpc

import (
	"encoding/json"
	"errors"
)

type responseMessage struct {
	Result json.RawMessage `json:"result"`
	Error  struct {
		Message string `json:"msg"`
	}
}

func (message responseMessage) result(result interface{}) error {
	if err := json.Unmarshal(message.Result, &result); err != nil {
		return err
	}
	return nil
}

func (message responseMessage) err() error {
	msg := message.Error.Message
	if msg == "" {
		return nil
	}
	return errors.New(msg)
}
