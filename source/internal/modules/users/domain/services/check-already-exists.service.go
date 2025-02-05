package services

import errs "github.com/Lucasvmarangoni/logella/err"

func (u *UserService) CheckAlreadyExists(email, cpf string) (bool, error) {

	// Terminar esse endpoint em outro momento. Não é importante agora.

	ok, err := u.UserRepository.CheckIfUserAlreadyExists(email, cpf, nil)
	if err != nil {
		return false, errs.Assertion(err).Context("u.Repository.CheckIfUserAlreadyExists")
	}
	return ok, nil
}
