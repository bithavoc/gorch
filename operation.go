package gorch

type Operation interface {
	Entry() *OperationEntry
	Close() error
}
