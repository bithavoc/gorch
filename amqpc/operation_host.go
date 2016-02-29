package amqpc

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/bithavoc/gorch"
	"github.com/streadway/amqp"
	"log"
)

type operationHost struct {
	host              *amqpcHost
	operation         *amqpcOperation
	entry             *gorch.OperationEntry
	shutdownSignal    chan struct{}
	requests          <-chan amqp.Delivery
	consumerName      string
	handler           gorch.OperationHandler
	operationShutdown chan<- *operationHost
}

func newOperationHost(host *amqpcHost, entry *gorch.OperationEntry) *operationHost {
	return &operationHost{
		host:           host,
		entry:          entry,
		shutdownSignal: make(chan struct{}, 2),
		consumerName:   uuid.NewRandom().String(),
	}
}

func (oph *operationHost) start(shutdown chan<- *operationHost) error {
	oph.operationShutdown = shutdown
	op, err := oph.host.cluster.operation(oph.entry.Name())
	if err != nil {
		return err
	}
	oph.operation = op
	oph.handler = oph.operation.Entry().Handler()
	const autoAck = false
	const noLocal = false
	const exclusive = false
	const noWait = false
	requests, err := oph.operation.channel.Consume(oph.operation.queueName(), oph.consumerName, autoAck, exclusive, noLocal, noWait, amqp.Table{})
	if err != nil {
		log.Printf("Host consume error: %s", err.Error())
		return err
	}
	oph.requests = requests
	go oph.loop()
	return nil
}

func (oph *operationHost) processRequest(request amqp.Delivery) {
	invocation := amqpcInvocation{}
	if err := json.Unmarshal(request.Body, &invocation); err != nil {
		request.Reject(false)
		log.Printf("Error deserializing operation request: %s", err.Error())
		return
	}
	result, err := oph.handler(invocation)
	if err != nil {
		log.Printf("Error executing handler for operation %s: %s", oph.operation.Entry().Name(), err.Error())
		request.Reject(false)
		return
	}
	if request.ReplyTo != "" {
		response := &responsePayload{}
		response.setError(err)
		response.Result = result

		responseData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error serializing result for operation %s: %s", oph.operation.Entry().Name(), err.Error())
			request.Reject(false)
			return
		}

		msg := amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         responseData,
		}
		if err := oph.operation.channel.Publish("", request.ReplyTo, false, false, msg); err != nil {
			log.Printf("Error Reply-To for operation %s: %s", oph.operation.Entry().Name(), err.Error())
			request.Reject(false)
			return
		}
	}

	request.Ack(false)
}

func (oph *operationHost) loop() {
	defer oph.operation.Close()
	defer oph.operation.channel.Cancel(oph.consumerName, false)
	func() {
		for {
			select {
			case request, ok := <-oph.requests:
				{
					if !ok {
						log.Printf("operation request transport has been closed")
						return
					}
					oph.processRequest(request)
				}
			case <-oph.shutdownSignal:
				{
					log.Printf("op host shutdown")
					return
				}
			}
		}
	}()
	log.Printf("Operation loop ended")
	oph.operationShutdown <- oph
}

func (oph *operationHost) shutdown() {
	if oph == nil {
		return
	}
	oph.shutdownSignal <- struct{}{}
}
