package amqpc

import (
	"fmt"
	"github.com/bithavoc/gorch"
	"github.com/streadway/amqp"
	"log"
)

type amqpcCluster struct {
	config     amqpcConfig
	connection *amqp.Connection
	registry   *gorch.OperationsRegistry
}

func newAmqpcCluster(config amqpcConfig) (*amqpcCluster, error) {
	cluster := &amqpcCluster{
		config:   config,
		registry: &gorch.OperationsRegistry{},
	}
	if err := cluster.open(); err != nil {
		return nil, err
	}
	return cluster, nil
}

func (cluster *amqpcCluster) open() error {
	if cluster.config.tlsConfig != nil {
		return cluster.dialTLS()
	}
	return cluster.dial()
}

func (cluster *amqpcCluster) dial() error {
	var err error
	cluster.connection, err = amqp.Dial(cluster.config.url)
	return err
}

func (cluster *amqpcCluster) dialTLS() error {
	var err error
	cluster.connection, err = amqp.DialTLS(cluster.config.url, cluster.config.tlsConfig)
	return err
}

func (cluster *amqpcCluster) Close() {
	if cluster == nil {
		return
	}
	err := cluster.connection.Close()
	if err != nil {
		log.Printf("Failed to close")
	}
}

func (cluster *amqpcCluster) Config() gorch.Config {
	return cluster.config
}

func (cluster *amqpcCluster) Operation(name string) (gorch.Operation, error) {
	return cluster.operation(name)
}

func (cluster *amqpcCluster) operation(name string) (*amqpcOperation, error) {
	entry := cluster.registry.Entry(name)
	if entry == nil {
		return nil, fmt.Errorf("Operation '%s' not found in registry", name)
	}
	op := newAmqpcOperation(cluster, entry)
	if err := op.open(); err != nil {
		return nil, err
	}
	return op, nil
}

func (cluster *amqpcCluster) Include(registry *gorch.OperationsRegistry) error {
	cluster.registry.Merge(registry)
	for _, entry := range registry.Operations() {
		if err := func() error {
			op, err := cluster.operation(entry.Name())
			if err != nil {
				return err
			}
			defer op.Close()
			if err := op.ensureTopology(); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

func (cluster *amqpcCluster) createChannel() (*amqp.Channel, error) {
	channel, err := cluster.connection.Channel()
	return channel, err
}
