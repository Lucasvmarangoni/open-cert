package services

import (
	"certifications/config"
	"certifications/internal/modules/users/domain/entities"
	"certifications/internal/modules/users/infra/repositories"
	"certifications/pkg/cache"
	"certifications/pkg/queue"
	"certifications/pkg/security"
	"encoding/json"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type UserService struct {
	Repository     repositories.UserRepository
	RabbitMQ       queue.IRabbitMQ
	MessageChannel chan amqp.Delivery
	ReturnChannel  chan error
	Memcached      cache.Mencacher[*entities.User]
	cacheUser      *entities.User
	aesKey         string
}

func NewUserService(
	repo repositories.UserRepository,
	rabbitMQ queue.IRabbitMQ,
	messageChannel chan amqp.Delivery,
	returnChannel chan error,
	memcached cache.Mencacher[*entities.User],
) *UserService {
	UserService := &UserService{
		Repository:     repo,
		RabbitMQ:       rabbitMQ,
		MessageChannel: messageChannel,
		ReturnChannel:  returnChannel,
		Memcached:      memcached,
		aesKey:         config.ReadSecretString(config.GetEnvString("security", "aes_key")),
	}
	return UserService
}

func (u *UserService) setToMemcacheIfNotNil(user *entities.User) {
	u.encrypt(user)
	if user != nil {
		u.Memcached.Set(user.ID.String(), user)
	}
}

func (u *UserService) encrypt(user *entities.User) error {
	var err error
	user.LastName, err = security.Encrypt(user.LastName, u.aesKey)
	if err != nil {
		return errs.Ctx(err, "security.Encrypt LastName")
	}
	user.Email, err = security.Encrypt(user.Email, u.aesKey)
	if err != nil {
		return errs.Ctx(err, "security.Encrypt Email")
	}
	user.CPF, err = security.Encrypt(user.CPF, u.aesKey)
	if err != nil {
		return errs.Ctx(err, "security.Encrypt CPF")
	}
	return nil
}

func (u *UserService) returnCachedUserIfExists() (*entities.User, bool) {
	if u.cacheUser != nil {
		u.decrypt(u.cacheUser)
		return u.cacheUser, true
	}
	return nil, false
}

func (u *UserService) getCache(key string) *entities.User {
	item, err := u.Memcached.Get(key)
	if err == nil && item != nil {
		var cachedUser entities.User
		err = json.Unmarshal(item.Value, &cachedUser)
		if err == nil {
			return &cachedUser
		} else {
			log.Error().Err(errs.Ctx(err, "json.Unmarshal")).Msg("Get User cache error")
		}
	}
	return nil
}
	

func (u *UserService) decrypt(user *entities.User) error {

	var err error

	user.LastName, err = security.Decrypt(user.LastName, u.aesKey)
	if err != nil {
		return errs.Ctx(err, "security.Decrypt LastName")
	}
	user.Email, err = security.Decrypt(user.Email, u.aesKey)
	if err != nil {
		return errs.Ctx(err, "security.Decrypt Email")
	}
	user.CPF, err = security.Decrypt(user.CPF, u.aesKey)
	if err != nil {
		return errs.Ctx(err, "security.Decrypt CPF")
	}
	return nil
}

func (u *UserService) deleteUserCache(id string) error {
	if u.cacheUser != nil {
		err := u.Memcached.Delete(id)
		if err != nil {
			return errs.Ctx(err, "u.Memcached.Delete")
		}
	}
	return nil
}
