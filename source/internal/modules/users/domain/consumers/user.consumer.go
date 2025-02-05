package consumers

import (
	"certifications/internal/modules/users/domain/entities"
	"certifications/internal/modules/users/infra/repositories"
	"certifications/pkg/queue"
	"context"
	"encoding/json"
	"net/http"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type UsersConsumeMessages struct {
	ctx        context.Context
	Repository *repositories.UserRepositoryDb
	RabbitMQ   *queue.RabbitMQ
}

func NewConsumeMessages(ctx context.Context, repository *repositories.UserRepositoryDb, rabbitMQ *queue.RabbitMQ) *UsersConsumeMessages {
	return &UsersConsumeMessages{
		ctx:        ctx,
		Repository: repository,
		RabbitMQ:   rabbitMQ,
	}
}

func (m *UsersConsumeMessages) UserConsumer(messageChannel chan amqp.Delivery, returnChannel chan error) {

	for message := range messageChannel {
		var user entities.User
		err := json.Unmarshal(message.Body, &user)
		if err != nil {
			returnChannel <- errs.Wrap(err, "json.Unmarshal", http.StatusInternalServerError)
			continue
		}

		err = m.Repository.Insert(&user, context.Background())
		if err != nil {
			returnChannel <- errs.Assertion(err).Context("Repository.Insert")
			continue
		} else {
			returnChannel <- nil
			log.Info().Str("context", "UserHandler").Msgf("User created successfully (%s)", user.ID)
		}
	}
}
