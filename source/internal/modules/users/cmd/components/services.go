package components

import (
	"certifications/internal/modules/users/domain/services"
	"certifications/pkg/queue"

	"github.com/streadway/amqp"
)

type ServiceFactory struct {
	UserService      *services.UserService
	CertifierService *services.CertifierService
}

func NewServiceFactory(
	repositories *UsersRespositoryFactory,
	memcached *UserCacheFactory,
	rabbitMQ queue.IRabbitMQ,
	messageChannel chan amqp.Delivery,
	returnChannel chan error,
) *ServiceFactory {
	return &ServiceFactory{
		UserService:      services.NewUserService(repositories.UserRepository, rabbitMQ, messageChannel, returnChannel, memcached.UserCache),
		CertifierService: services.NewCertifierService(repositories.CertifierRepository, memcached.CertifierCache),
	}
}
