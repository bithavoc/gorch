package amqpc

import (
	"github.com/bithavoc/gorch"
)

type amqpcConfig struct {
	url string
}

func NewAmqpcClusterConfig(url string) gorch.Config {
	return amqpcConfig{
		url: url,
	}
}

func (config amqpcConfig) Connect() (gorch.Cluster, error) {
	cluster, err := newAmqpcCluster(config)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (config amqpcConfig) CreateHost() gorch.Host {
	return nil
}
