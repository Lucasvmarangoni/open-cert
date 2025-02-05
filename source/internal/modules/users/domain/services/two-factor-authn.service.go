package services

import (
	"context"
	"fmt"
	"net/http"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/pquerna/otp/totp"
)

type otpResponseOutput struct {
	Base32      string
	Otpauth_url string
	AccessToken string
}

func (u *UserService) GenerateTOTP(id string) (*otpResponseOutput, error) {
	// deslogar o usuário após ele habilitar o 2fa??
	var otpResponse *otpResponseOutput

	user, err := u.FindById(id, context.Background())
	if err != nil {
		return nil, errs.Assertion(err).Context("u.FindById")
	}

	if user.OtpEnabled {
		return nil, errs.Wrap(err, "Already enabled 2FA", http.StatusBadRequest)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "certificator",
		AccountName: id,
		SecretSize:  15,
	})
	if err != nil {
		return nil, errs.Wrap(err, "totp.Generate", http.StatusInternalServerError)
	}
	tokenString, err := u.generateTempJWT(user)
	if err != nil {
		return nil, errs.Assertion(err).Context("u.generateJWT")
	}
	otpResponse = &otpResponseOutput{
		Base32:      key.Secret(),
		Otpauth_url: key.URL(),
		AccessToken: tokenString,
	}
	user.OtpSecret = otpResponse.Base32
	user.OtpAuthUrl = otpResponse.Otpauth_url

	err = u.UserRepository.UpdateOTP(user, context.Background())
	if err != nil {
		return nil, errs.Assertion(err).Context("u.Repository.UpdateTOTP")
	}

	err = u.deleteUserCache(id)
	if err != nil {
		return nil, errs.Assertion(err).Context("u.deleteUserCache")
	}
	u.setToMemcacheIfNotNil(*user)

	return otpResponse, nil
}

func (u *UserService) VerifyTOTP(id, tempToken string) (authData, error) {

	user, err := u.FindById(id, context.Background())
	if err != nil {
		return authData{}, errs.Assertion(err).Context("u.FindById")
	}

	valid := totp.Validate(tempToken, user.OtpSecret)
	if !valid {
		return authData{}, errs.Wrap(fmt.Errorf("token is invalid"), "u.FindById", http.StatusUnauthorized)
	}
	if !user.OtpEnabled {
		user.OtpEnabled = true
		err = u.UserRepository.UpdateOTPVerify(user, context.Background())
		if err != nil {
			return authData{}, errs.Assertion(err).Context("u.Repository.UpdateTOTPVerify")
		}
		err = u.deleteUserCache(id)
		if err != nil {
			return authData{}, errs.Assertion(err).Context("u.deleteUserCache")
		}
		u.setToMemcacheIfNotNil(*user)
	}
	tokenString, err := u.generateJWT(user)
	if err != nil {
		errs.Assertion(err).Context("u.generateJWT")
	}

	return authData{
		Token:            tokenString,
		IsTemporaryToken: user.OtpEnabled,
	}, nil
}

func (u *UserService) DisableOTP(id string) error {
	user, err := u.FindById(id, context.Background())
	if err != nil {
		return errs.Assertion(err).Context("u.FindById")
	}
	user.OtpEnabled = false

	err = u.UserRepository.UpdateOTPVerify(user, context.Background())
	if err != nil {
		return errs.Assertion(err).Context("u.Repository.UpdateTOTPVerify")
	}
	err = u.deleteUserCache(id)
	if err != nil {
		return errs.Assertion(err).Context("u.deleteUserCache")
	}
	u.setToMemcacheIfNotNil(*user)

	return nil
}
