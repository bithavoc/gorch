package amqpc

import (
	"crypto/tls"
	"github.com/bithavoc/gorch"
)

type amqpcConfig struct {
	url       string
	tlsConfig *tls.Config
}

func NewAmqpcClusterConfig(url string, tlsConfig *tls.Config) gorch.Config {
	return amqpcConfig{
		url:       url,
		tlsConfig: tlsConfig,
	}
}

func (config amqpcConfig) Connect() (gorch.Cluster, error) {
	return config.connect()
}

func (config amqpcConfig) connect() (*amqpcCluster, error) {
	cluster, err := newAmqpcCluster(config)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (config amqpcConfig) Host() (gorch.Host, error) {
	return newAmqpcHost(config)
}
