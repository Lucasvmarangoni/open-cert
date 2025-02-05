package services

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/Lucasvmarangoni/logella/err"
)

func (u *UserService) Delete(id string) error {

	user, err := u.FindById(id, context.Background())
	if err != nil {
		return errs.Assertion(err).Context("u.Repository.Delete")
	}
	if user.Certifier {
		return errs.Wrap(errors.New(
			"you have a certifier account. it's necessary request the deletion of the Certifier account first"),
			"u.deleteCache",
			http.StatusForbidden,
		)
	}
	err = u.UserRepository.Delete(id, context.Background())
	if err != nil {
		return errs.Assertion(err).Context("u.Repository.Delete")
	}
	err = u.deleteUserCache(id)
	if err != nil {
		return errs.Assertion(err).Context("u.deleteCache")
	}
	return nil
}
