package router

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/the-psyducks/metrics-service/src/config"
	"github.com/the-psyducks/metrics-service/src/models"
	"github.com/the-psyducks/metrics-service/src/repository"
	"github.com/the-psyducks/metrics-service/src/service"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const NewUserType = "NEW_USER"
const NewRegistryType = "NEW_REGISTRY"
const UserBlockedType = "USER_BLOCKED"
const LoginAttemptType = "LOGIN_ATTEMPT"

type AmpqRouter struct {
	messagesChan   <-chan amqp.Delivery
	metricsService *service.MetricsService
}

func NewRabbitRouter() (*AmpqRouter, error) {
	queueName := os.Getenv("CLOUDAMQP_QUEUE")
	queueUrl := os.Getenv("CLOUDAMQP_URL")
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	conn, err := amqp.Dial(queueUrl)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	messagesChan, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	metricsDb, err := repository.CreateMetricsDatabases(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics database: %w", err)
	}
	metricsService := service.NewMetricsService(metricsDb)

	return &AmpqRouter{
		messagesChan:   messagesChan,
		metricsService: metricsService,
	}, nil
}

func (r *AmpqRouter) Run() {
	go func() {
		for message := range r.messagesChan {
			var queueMessage models.QueueMessage
			err := json.Unmarshal(message.Body, &queueMessage)
			if err != nil {
				slog.Warn(fmt.Sprintf("error unmarshalling message: %v", err))
			}

			switch queueMessage.Type {
			case LoginAttemptType:
				var event models.QueueLoginAttempt
				if err := json.Unmarshal(message.Body, &event); err != nil {
					slog.Warn(fmt.Sprintf("error unmarshalling login attempt message: %v", err))
				}
				slog.Info(fmt.Sprintf("received login attempt: %v", event.Message))
				if err := r.metricsService.RecordLoginAttempt(event.Message); err != nil {
					slog.Error(fmt.Sprintf("error recording login attempt: %v", err))
					return
				}
			case UserBlockedType:
				fallthrough
			case NewRegistryType:
				fallthrough
			case NewUserType:
				slog.Warn(fmt.Sprintf("message type %s not implemented", queueMessage.Type))
			default:
				slog.Warn(fmt.Sprintf("message type %s not recognized", queueMessage.Type))
			}
		}
	}()
}
