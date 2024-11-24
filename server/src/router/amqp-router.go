package router

import (
	"fmt"
	"os"
	amqp "github.com/rabbitmq/amqp091-go"
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
	queueName := os.Getenv("CLOUDAMQP_QUEUE")
	queueUrl := os.Getenv("CLOUDAMQP_URL")
	fmt.Println("dialing: ", queueUrl)
	conn, err := amqp.Dial(queueUrl)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	fmt.Println("declaring queue: ", queueName)
	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("consuming queue: ", queueName)
	messagesChan, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	
	fmt.Println("consumed")

	return &AmpqRouter{messagesChan: messagesChan}, nil
}
