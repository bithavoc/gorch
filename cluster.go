package gorch

type Cluster interface {
	Config() Config
	Operation(name string) (Operation, error)
	Close()
	Include(registry *OperationsRegistry) error
}
