package services

import (
	"certifications/internal/modules/users/domain/entities"
	"time"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/go-chi/jwtauth"
)

func (u *UserService) Authn(email, password string, jwt *jwtauth.JWTAuth, jwtExpiresIn int) (string, error) {
	var user *entities.User
	var err error	

	user, err = u.FindByEmail(email, nil)

	if err != nil {
		return "", errs.Ctx(err, "u.FindByEmail")
	}

	err = user.ValidateHashPassword(password)
	if err != nil {
		return "", errs.Ctx(err, "user.ValidateHashPassword")
	}
	user.Password = ""

	tokenString, err := u.generateJWT(user, jwt, jwtExpiresIn)
	if err != nil {
		errs.Ctx(err, "u.generateJWT")
	}
	return tokenString, nil
}

func (u *UserService) generateJWT(user *entities.User, jwt *jwtauth.JWTAuth, jwtExpiresIn int) (string, error) {
	_, tokenString, err := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})
	if err != nil {
		errs.Ctx(err, "jwt.Encode")
	}
	return tokenString, nil
}
