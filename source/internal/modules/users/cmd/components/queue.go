package components

import (
	"certifications/config"
	"certifications/internal/modules/users/domain/consumers"
	"certifications/internal/modules/users/infra/repositories"
	"certifications/pkg/queue"
	"context"

	"github.com/streadway/amqp"
)

func Queues(ctx context.Context, userRepository *repositories.UserRepositoryDb) (chan amqp.Delivery, chan error, *queue.RabbitMQ, *amqp.Channel) {
	rabbitMQ := queue.NewRabbitMQ()
	userConsumeMessages := consumers.NewConsumeMessages(ctx, userRepository, rabbitMQ)
	ch := rabbitMQ.Connect()

	messageChannel, returnChannel := userConsumer(userConsumeMessages, rabbitMQ)

	return messageChannel, returnChannel, rabbitMQ, ch
}

func userConsumer(
	userConsumer *consumers.UsersConsumeMessages,
	rabbitMQ *queue.RabbitMQ,
) (chan amqp.Delivery, chan error) {

	messageChannel := make(chan amqp.Delivery)
	returnChannel := make(chan error)
	rabbitMQ.Consume(messageChannel, config.GetEnvString("rabbitMQ", "routingkey_userCreate"))

	go userConsumer.UserConsumer(messageChannel, returnChannel)

	return messageChannel, returnChannel
}
