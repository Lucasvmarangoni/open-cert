package services

import (
	"context"

	errs "github.com/Lucasvmarangoni/logella/err"
)

func (u *UserService) Update(id, password, newEmail, newPassword string) error {

	user, err := u.FindByIdWiyhPassword(id, nil)
	if err != nil {
		return errs.Assertion(err).Context("u.FindById")
	}
	err = user.ValidateHashPassword(password)
	if err != nil {
		return errs.Assertion(err).Context("user.ValidateHashPassword")
	}
	err = u.deleteUserCache(id)
	if err != nil {
		return errs.Assertion(err).Context("u.deleteCache")
	}
	if newPassword != "" {
		password = newPassword
	}
	err = user.Update(newEmail, password)
	if err != nil {
		return errs.Assertion(err).Context("user.Update")
	}
	err = u.UserRepository.Update(user, context.Background())
	if err != nil {
		return errs.Assertion(err).Context("Repository.Update")
	}
	u.setToMemcacheIfNotNil(*user)
	return nil
}
