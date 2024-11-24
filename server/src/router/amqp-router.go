package router

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AmpqRouter struct {
	messagesChan <-chan amqp.Delivery
	AmqpChannel	*amqp.Channel
	AmqpConn	*amqp.Connection
}

func (r *AmpqRouter) Run() {
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case msg := <-r.messagesChan:
				fmt.Println("Message: ", string(msg.Body))
			case <-sigchan:
				fmt.Println("Received termination signal")

				if err := r.AmqpChannel.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to close AMQP channel: %v\n", err)
				}
				if err := r.AmqpConn.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to close AMQP connection: %v\n", err)
				}
				return
			}
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

	return &AmpqRouter{
		messagesChan: messagesChan,
		AmqpChannel: channel,
		AmqpConn: conn,
		}, nil
}
