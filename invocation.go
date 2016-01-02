package gorch

type Invocation interface {
	Input(input interface{}) error
}
