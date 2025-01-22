package components

import (
	"certifications/internal/modules/users/domain/entities"
	"certifications/internal/modules/users/domain/services"
	"certifications/internal/modules/users/infra/repositories"
	"certifications/pkg/cache"
	"certifications/pkg/queue"

	"github.com/streadway/amqp"
)

type ServiceFactory struct {
	UserService *services.UserService
}

func NewServiceFactory[T any](
	userRepository *repositories.UserRepositoryDb,
	rabbitMQ queue.IRabbitMQ,
	messageChannel chan amqp.Delivery,
	returnChannel chan error,
	memcached cache.Mencacher[*entities.User],
) *ServiceFactory {
	return &ServiceFactory{
		UserService: services.NewUserService(userRepository, rabbitMQ, messageChannel, returnChannel, memcached),
	}
}
