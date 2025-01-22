package services

import (
	"certifications/config"
	"certifications/internal/modules/users/domain/entities"
	"encoding/json"
	"sync"

	"github.com/Lucasvmarangoni/logella/err"
)

func (u *UserService) Create(name, lastName, cpf, email, password string) error {

	newUser, err := entities.NewUser(name, lastName, cpf, email, password)
	if err != nil {
		return errs.Ctx(err, "entities.NewUser")
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		return errs.Ctx(err, "json.Marshal")
	}

	err = u.RabbitMQ.Publish(string(userJSON), "application/json", config.GetEnvString("rabbitMQ", "exchange"), config.GetEnvString("rabbitMQ", "queue_user"), config.GetEnvString("rabbitMQ", "routingkey_userCreate"))
	if err != nil {
		return errs.Ctx(err, "RabbitMQ.Publish")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = <-u.ReturnChannel
	}()
	wg.Wait()
	if err != nil {
		return errs.Ctx(err, "UserConsumer")
	}

	u.setToMemcacheIfNotNil(newUser)

	return nil
}
