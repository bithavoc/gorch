package amqpc

type responsePayload struct {
	Result interface{} `json:"result"`
	Error  struct {
		Message string `json:"msg"`
	}
}

func (payload *responsePayload) setError(err error) {
	if err == nil {
		return
	}
	payload.Error.Message = err.Error()
}
