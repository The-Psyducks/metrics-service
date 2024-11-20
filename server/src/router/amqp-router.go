package router

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"os/signal"
	"syscall"
)

type AmpqRouter struct {
	messagesChan <-chan amqp.Delivery
}

func (r *AmpqRouter) Run() {
	go func() {
		for message := range r.messagesChan {
			println(string(message.Body))
		}
	}()
}

func NewRabbitRouter() (*AmpqRouter, error) {
	queueName := os.Getenv("RABBIT_QUEUE")
	queueUrl := os.Getenv("RABBIT_URL")
	conn, err := amqp.Dial(queueUrl)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	messagesChan, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	return &AmpqRouter{messagesChan: messagesChan}, nil
}
