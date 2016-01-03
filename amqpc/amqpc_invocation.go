package amqpc

import (
	"bytes"
	"encoding/json"
)

type amqpcInvocation struct {
	ModeData  invocationMode  `json:"mode"`
	InputData json.RawMessage `json:"input"`
}

func (req amqpcInvocation) Input(input interface{}) error {
	buff := bytes.NewBuffer(req.InputData)
	decoder := json.NewDecoder(buff)
	decoder.UseNumber()
	err := decoder.Decode(&input)
	return err
}
