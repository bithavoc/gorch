package amqpc

type invocationMode string

const invocationModeCall invocationMode = "call"

type invocationPayload struct {
	Mode  invocationMode `json:"mode"`
	Input interface{}    `json:"input"`
}
