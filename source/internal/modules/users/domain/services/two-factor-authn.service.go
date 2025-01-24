package services

import (
	"context"
	"fmt"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/pquerna/otp/totp"
)

type otpResponseOutput struct {
	Base32      string
	Otpauth_url string
}

func (u *UserService) GenerateTOTP(id string) (*otpResponseOutput, error) {

	var otpResponse *otpResponseOutput

	user, err := u.FindById(id, context.Background())
	if err != nil {
		return nil, errs.Ctx(err, "u.FindById")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "certificator",
		AccountName: id,
		SecretSize:  15,
	})
	if err != nil {
		return nil, errs.Ctx(err, "totp.Generate")
	}

	err = u.Repository.UpdateOTP(user, context.Background())
	if err != nil {
		return nil, errs.Ctx(err, "u.Repository.UpdateOTP")
	}

	otpResponse = &otpResponseOutput{
		Base32:      key.Secret(),
		Otpauth_url: key.URL(),
	}
	return otpResponse, nil
}

func (u *UserService) VerifyTOTP(id, token, isValidate string) error {

	user, err := u.FindById(id, context.Background())
	if err != nil {
		return errs.Ctx(err, "u.FindById")
	}

	valid := totp.Validate(token, user.OtpSecret)
	if !valid {
		return errs.Ctx(fmt.Errorf("token is invalid"), "u.FindById")
	}

	if isValidate == "1" {
		user.OtpEnabled = true
	}

	err = u.Repository.UpdateOTPVerify(user, context.Background())
	if err != nil {
		return errs.Ctx(err, "u.Repository.UpdateOTPVerify")
	}
	return nil
}

func (u *UserService) DisableOTP(id string) error {
	user, err := u.FindById(id, context.Background())
	if err != nil {
		return errs.Ctx(err, "u.FindById")
	}

	user.OtpEnabled = false

	err = u.Repository.UpdateOTPVerify(user, context.Background())
	if err != nil {
		return errs.Ctx(err, "u.Repository.UpdateOTPVerify")
	}
	return nil
}
