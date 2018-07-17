package carrot

import (
	"github.com/streadway/amqp"
)

//Subscriber is a consumer component to Rabbit
type Subscriber struct {
	client *BrokerClient
}

//MessageContext manager received message from rabbit and ack process
type MessageContext struct {
	Message    Message
	delivery   amqp.Delivery
	subscriber *Subscriber
}

//Ack message to server
func (ctx *MessageContext) Ack() error {
	return ctx.delivery.Ack(false)
}

//Nack message to server if requeue = true the message will be sent to same queue
func (ctx *MessageContext) Nack(requeue bool) error {
	return ctx.delivery.Nack(false, requeue)
}

//RedirectTo redirect message to other exchange
func (ctx *MessageContext) RedirectTo(exchange, routingKey string) error {
	ch, err := ctx.subscriber.client.Channel()
	if err == nil {
		err = ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
			Body:            ctx.Message.Data,
			ContentEncoding: ctx.Message.Encoding,
			ContentType:     ctx.Message.ContentType,
			Headers:         ctx.Message.Headers,
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		})
	}
	return err
}

//SubscribeWorker is the worker handler for queues
type SubscribeWorker struct {
	Scale   uint
	Handler func(*MessageContext) error
	AutoAck bool
	Queue   string
}

//Subscribe binds a worker to queue on Rabbit
func (sub *Subscriber) Subscribe(worker SubscribeWorker) error {
	var i uint = 0
	messageHandler := func(worker SubscribeWorker, msgsChan <-chan amqp.Delivery) {
		for message := range msgsChan {
			context := new(MessageContext)
			context.Message = Message{
				ContentType: message.ContentType,
				Data:        message.Body,
				Encoding:    message.ContentEncoding,
			}
			context.delivery = message
			context.subscriber = sub
			worker.Handler(context)
		}
	}
	for i = 0; i < worker.Scale; i++ {
		ch, err := sub.client.Channel()
		if err == nil {
			msgs, err := ch.Consume(
				worker.Queue,   // queue
				"",             // consumer
				worker.AutoAck, // auto-ack
				false,          // exclusive
				false,          // no-local
				false,          // no-wait
				nil,            // args
			)
			if err == nil {
				go messageHandler(worker, msgs)
			}
		} else {
			sub.client.channel = nil
			i--
		}
	}
	return nil
}

//NewSubscriber creates a new Subscriber for Rabbit
func NewSubscriber(client *BrokerClient) *Subscriber {
	subs := new(Subscriber)
	subs.client = client
	return subs
}
