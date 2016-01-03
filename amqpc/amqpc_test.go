package amqpc

import (
	"github.com/bithavoc/gorch"
	"log"
	"testing"
)

func TestAmqp(t *testing.T) {
	const opName = "clusterTopologyTestOperation"
	config := NewAmqpcClusterConfig(amqpcServerTestingURL)
	cluster, err := config.Connect()
	if err != nil {
		t.Fatalf("Amqpc cluster connection error: %s", err.Error())
	}

	host, err := config.Host()
	if err != nil {
		t.Fatalf("Amqpc host error: %s", err.Error())
	}
	defer host.Shutdown()

	done := make(chan error, 2)
	go func(done chan error) {
		hostRegistry := &gorch.OperationsRegistry{}
		hostRegistry.Register(opName).Host(func(invocation gorch.Invocation) (interface{}, error) {
			var args struct {
				Greeting string `json:"greeting"`
			}
			err := invocation.Input(&args)
			if err != nil {
				return nil, err
			}
			return struct {
				Response string `json:"response"`
			}{
				Response: args.Greeting + " back",
			}, nil
		})
		if err := host.Mount(hostRegistry); err != nil {
			t.Fatalf("Amqpc host registry error: %s", err.Error())
		}
		log.Printf("Serving")
		if err := host.Serve(); err != nil {
			t.Fatalf("Amqpc host serve error: %s", err.Error())
		}
		log.Printf("Done serving test host")
		done <- nil
	}(done)

	registry := &gorch.OperationsRegistry{}
	registry.Register(opName)
	if err := cluster.Include(registry); err != nil {
		t.Fatalf("Amqpc cluster register operations error: %s", err.Error())
	}
	op, err := cluster.Operation(opName)
	if err != nil {
		t.Fatalf("Amqpc cluster operation failed: %s", err.Error())
	}
	defer op.Close()

	args := struct {
		Greeting string `json:"greeting"`
	}{
		Greeting: "hola",
	}
	out := struct {
		Response string `json:"response"`
	}{}
	err = op.Call(args).Resulting(&out)
	if err != nil {
		t.Fatalf("Amqpc cluster operation failed: %s", err.Error())
	}
	if out.Response != "hola back" {
		t.Fatalf("Response was not received")
	}
}
