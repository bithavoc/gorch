package amqpc

import (
	"errors"
	"fmt"
	"github.com/bithavoc/gorch"
	"log"
)

type amqpcHost struct {
	config         amqpcConfig
	cluster        *amqpcCluster
	shutdown       chan struct{}
	serving        bool
	operationHosts map[string]*operationHost
}

func newAmqpcHost(config amqpcConfig) (*amqpcHost, error) {
	cluster, err := config.connect()
	if err != nil {
		return nil, err
	}
	host := &amqpcHost{
		config:         config,
		cluster:        cluster,
		shutdown:       make(chan struct{}, 2),
		operationHosts: make(map[string]*operationHost),
	}
	return host, nil
}

func (host *amqpcHost) Mount(registry *gorch.OperationsRegistry) error {
	if host.serving {
		return errors.New("Unable to mount, Host is already serving")
	}
	analysis := registry.Analyze()
	if !analysis.IsFullyHosted() {
		firstMissingHost := analysis.MissingHosts[0]
		err := fmt.Errorf("Unable to mount registry, missing host for operation '%s'", firstMissingHost.Name())
		for _, op := range analysis.MissingHosts {
			log.Printf("Missing host for operation '%s'", op.Name())
		}
		return err
	}
	host.cluster.Include(registry)
	return nil
}

func (host *amqpcHost) isServing(name string) bool {
	_, ok := host.operationHosts[name]
	return ok
}

func (host *amqpcHost) add(oph *operationHost) {
	host.operationHosts[oph.operation.Entry().Name()] = oph
}

func (host *amqpcHost) Serve() error {
	if host.serving {
		return errors.New("Unable to host, Host is already serving")
	}
	host.serving = true
	defer host.cluster.Close()
	for _, entry := range host.cluster.registry.Operations() {
		log.Printf("Starting operation host for entry: %s", entry.Name())
		if host.isServing(entry.Name()) {
			continue
		}
		oph := newOperationHost(host, entry)
		defer oph.shutdown()
		if err := oph.start(); err != nil {
			return err
		}
		host.add(oph)
	}
	<-host.shutdown
	return nil
}

func (host *amqpcHost) Shutdown() error {
	if host == nil {
		return nil
	}
	host.shutdown <- struct{}{}
	return nil
}
