package gorch

type Operation interface {
	Entry() *OperationEntry
	Close() error
	Call(input interface{}) ProcedureCall
}
