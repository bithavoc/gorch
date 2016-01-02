package amqpc

import (
	"github.com/bithavoc/gorch"
	"testing"
)

func TestAmqpOperationEmitting(t *testing.T) {
	config := NewAmqpcClusterConfig(amqpcServerTestingURL)
	cluster, err := config.Connect()
	if err != nil {
		t.Fatalf("Amqpc cluster connection error: %s", err.Error())
	}
	registry := &gorch.OperationsRegistry{}
	opName := "clusterTopologyTestOperation"
	registry.Register(opName)
	if err := cluster.Include(registry); err != nil {
		t.Fatalf("Amqpc cluster register operations error: %s", err.Error())
	}
	op, err := cluster.Operation(opName)
	if err != nil {
		t.Fatalf("Amqpc cluster operation failed: %s", err.Error())
	}
	defer op.Close()
}
