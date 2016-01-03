package gorch

type Host interface {
	Mount(registry *OperationsRegistry) error
	Serve() error
	Shutdown() error
}
