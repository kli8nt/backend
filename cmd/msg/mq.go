package msg

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQConfig struct {
	Host string
	Port string
	User string
	Pass string
}

type MQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

type Queue struct {
	name string
	q    amqp.Queue
	ch   *amqp.Channel
}

func onError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (mq *MQ) Init(config MQConfig) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.User, config.Pass, config.Host, config.Port)
	conn, err := amqp.Dial(uri)
	(*mq).connection = conn
	onError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	(*mq).channel = ch
	onError(err, "Failed to open a channel")

	fmt.Println("Successfully connected to RabbitMQ")
}

func (mq *MQ) Queue(name string) Queue {
	q, err := (*mq).channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	onError(err, "Failed to declare a queue")

	queue := Queue{
		name: name,
		q:    q,
		ch:   (*mq).channel,
	}
	return queue
}

func (q *Queue) Consume(cb func(msg []byte)) {
	msgs, err := (*q).ch.Consume(
		(*q).name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	onError(err, "Failed to register a consumer")
	log.Printf(" [*] Waiting for messages. for Queue %s", (*q).name)
	for d := range msgs {
		cb(d.Body)
	}
}

func (q *Queue) Publish(msg []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := (*q).ch.PublishWithContext(ctx,
		"",        // exchange
		(*q).name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})

	onError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", msg)
	defer cancel()
}

func (mq *MQ) Close() {
	(*mq).connection.Close()
	(*mq).channel.Close()
}
