package gorch

type Config interface {
	Connect() (Cluster, error)
	Host() (Host, error)
}
