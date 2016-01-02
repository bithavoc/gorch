package gorch

type Config interface {
	Connect() (Cluster, error)
	CreateHost() Host
}
