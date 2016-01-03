package gorch

type OperationHandler func(invocation Invocation) (interface{}, error)
