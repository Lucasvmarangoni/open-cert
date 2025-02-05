package services

import (
	"certifications/config"
	"certifications/internal/modules/users/domain/entities"
	"net/http"
	"time"

	errs "github.com/Lucasvmarangoni/logella/err"
)

type authData struct {
	Token            string
	IsTemporaryToken bool
}

func (u *UserService) Authn(email, password string) (authData, error) {
	var user *entities.User
	var err error

	user, err = u.FindByEmail(email, nil)
	if err != nil {
		return authData{}, errs.Assertion(err).Context("u.FindByEmail")
	}

	err = user.ValidateHashPassword(password)
	if err != nil {
		return authData{}, errs.Assertion(err).Context("user.ValidateHashPassword")
	}
	user.Password = ""

	if user.OtpEnabled {
		tokenString, err := u.generateTempJWT(user)
		if err != nil {
			errs.Assertion(err).Context("u.generateJWT")
		}
		return authData{
			Token:            tokenString,
			IsTemporaryToken: true,
		}, nil
	}

	tokenString, err := u.generateJWT(user)
	if err != nil {
		errs.Assertion(err).Context("u.generateJWT")
	}
	return authData{
		Token:            tokenString,
		IsTemporaryToken: false,
	}, nil
}

func (u *UserService) generateJWT(user *entities.User) (string, error) {
	tokenAuth := config.GetTokenAuth()
	jwtExpiresIn := config.GetEnvInt("jwt", "expiredIn")

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Minute * time.Duration(jwtExpiresIn)).Unix(),
	})
	if err != nil {
		return "", errs.Wrap(err, "jwt.Encode", http.StatusInternalServerError)
	}
	return tokenString, nil
}

func (u *UserService) generateTempJWT(user *entities.User) (string, error) {
	tokenAuthTemp := config.GetTokenAuthTemp()
	jwtExpiresInTemp := config.GetEnvInt("jwt", "expiredIn_temp")

	_, tokenString, err := tokenAuthTemp.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Minute * time.Duration(jwtExpiresInTemp)).Unix(),
	})
	if err != nil {
		return "", errs.Wrap(err, "jwt.Encode", http.StatusInternalServerError)
	}
	return tokenString, nil
}
