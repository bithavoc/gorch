package amqpc

import (
	"encoding/json"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/streadway/amqp"
	"log"
)

type amqpcProcedureCall struct {
	input     interface{}
	channel   *amqp.Channel
	operation *amqpcOperation
}

func newAmqpcProcedureCall(operation *amqpcOperation, input interface{}) *amqpcProcedureCall {
	call := &amqpcProcedureCall{
		operation: operation,
		channel:   operation.channel,
		input:     input,
	}

	return call
}

func (call *amqpcProcedureCall) Resulting(result interface{}) error {
	const durable = false
	const autoDelete = true
	const exclusive = true
	const noWait = false
	/*queue, err := call.channel.QueueDeclare("", durable, autoDelete, exclusive, noWait, amqp.Table{})
	if err != nil {
		log.Printf("Procedure: %s", err.Error())
		return err
	}
	defer call.channel.QueueDelete(queue.Name, false, false, false)*/
	consumerName := uuid.NewRandom().String()
	const autoAck = true
	const noLocal = false
	replies, err := call.channel.Consume("amq.rabbitmq.reply-to", consumerName, autoAck, exclusive, noLocal, noWait, amqp.Table{})
	if err != nil {
		log.Printf("Procedure consume error: %s", err.Error())
		return err
	}
	defer call.channel.Cancel(consumerName, false)
	payload := invocationPayload{
		Mode:  invocationModeCall,
		Input: call.input,
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to serialize input with error: %s", err.Error())
	}
	const mandatory = true
	const immediate = false
	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         payloadData,
		ReplyTo:      "amq.rabbitmq.reply-to",
	}
	if err := call.channel.Publish(call.operation.exchangeName(), call.operation.exchangeKey(), mandatory, immediate, msg); err != nil {
		return err
	}
	select {
	case reply := <-replies:
		{
			responseMsg := responseMessage{}
			if err := json.Unmarshal(reply.Body, &responseMsg); err != nil {
				return err
			}
			if err := responseMsg.result(&result); err != nil {
				return err
			}
		}
	}
	return nil
}
