package amqpc

import (
	"github.com/bithavoc/gorch"
	"github.com/streadway/amqp"
	"log"
)

type amqpcOperation struct {
	cluster *amqpcCluster
	entry   *gorch.OperationEntry
	channel *amqp.Channel
}

func newAmqpcOperation(cluster *amqpcCluster, entry *gorch.OperationEntry) *amqpcOperation {
	return &amqpcOperation{
		cluster: cluster,
		entry:   entry,
	}
}

func (op *amqpcOperation) open() (err error) {
	op.channel, err = op.cluster.connection.Channel()
	return
}

func (op *amqpcOperation) Close() error {
	if op == nil {
		return nil
	}
	if op.channel != nil {
		return op.channel.Close()
	}
	op.channel = nil
	return nil
}

func (op *amqpcOperation) exchangeName() string {
	return op.entry.Name()
}

func (op *amqpcOperation) exchangeKey() string {
	return op.exchangeName()
}

func (op *amqpcOperation) queueName() string {
	return op.exchangeName()
}

func (op *amqpcOperation) ensureTopology() error {
	const durable = true
	const autoDelete = false
	const internal = false
	const noWait = false
	if err := op.channel.ExchangeDeclare(op.exchangeName(), "fanout", durable, autoDelete, internal, noWait, amqp.Table{}); err != nil {
		log.Printf("Exchange: %s", err.Error())
		return err
	}
	const exclusive = false
	_, err := op.channel.QueueDeclare(op.queueName(), durable, autoDelete, exclusive, noWait, amqp.Table{})
	if err != nil {
		log.Printf("Queue declare error: %s", err.Error())
		return err
	}
	if err := op.channel.QueueBind(op.queueName(), op.exchangeKey(), op.exchangeName(), noWait, amqp.Table{}); err != nil {
		log.Printf("Queue bind error: %s", err.Error())
		return err
	}
	return nil
}

func (op *amqpcOperation) Entry() *gorch.OperationEntry {
	return op.entry
}
