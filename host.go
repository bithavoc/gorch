package gorch

type Host interface {
	Mount(registry *OperationsRegistry) error
	Serve() chan<- HostTermination
	Shutdown() error
}
