package services

import (
	"certifications/internal/modules/users/domain/entities"
	pkg_entities "certifications/pkg/entities"
	"context"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/rs/zerolog/log"
)

func (u *UserService) FindById(id string, ctx context.Context) (*entities.User, error) {

	u.cacheUser = u.getCache(id)
	if cachedUser, ok := u.returnCachedUserIfExists(); ok {
		log.Log().Msgf("cached user", cachedUser.Email)
		return cachedUser, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	parsedId, err := pkg_entities.ParseID(id)
	if err != nil {
		return nil, errs.Assertion(err).Context("pkg_entities.ParseID")
	}
	user, err := u.UserRepository.FindById(parsedId, ctx)
	if err != nil {
		return nil, errs.Assertion(err).Context("Repository.FindById")
	}
	u.setToMemcacheIfNotNil(*user)
	user.Password = ""
	return user, nil
}

func (u *UserService) FindByIdWiyhPassword(id string, ctx context.Context) (*entities.User, error) {

	u.cacheUser = u.getCache(id)
	if cachedUser, ok := u.returnCachedUserIfExists(); ok {
		log.Log().Msgf("cached user", cachedUser.Email)
		return cachedUser, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	parsedId, err := pkg_entities.ParseID(id)
	if err != nil {
		return nil, errs.Assertion(err).Context("pkg_entities.ParseID")
	}
	user, err := u.UserRepository.FindById(parsedId, ctx)
	if err != nil {
		return nil, errs.Assertion(err).Context("Repository.FindById")
	}
	u.setToMemcacheIfNotNil(*user)
	return user, nil
}

func (u *UserService) FindByEmail(email string, ctx context.Context) (*entities.User, error) {

	if cachedUser, ok := u.returnCachedUserIfExists(); ok {
		log.Log().Msgf("cached user", cachedUser.Email)
		return cachedUser, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	user, err := u.UserRepository.FindByEmail(email, ctx)
	if err != nil {
		return nil, errs.Assertion(err).Context("Repository.FindByEmail")
	}
	u.setToMemcacheIfNotNil(*user)	
	return user, nil
}
