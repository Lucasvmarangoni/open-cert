package dto

import "github.com/asaskevich/govalidator"

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type UserInput struct {
	Name     string `json:"name" valid:"required,matches(^[a-zA-Z ]+$),length(3|10)"`
	LastName string `json:"last_name" valid:"required,matches(^[a-zA-Z ]+$),length(3|50)"`
	CPF      string `json:"cpf" valid:"matches(^[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}-[0-9]{2}$)"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"required"`
}


type AuthenticationInput struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required"`
}

type GetJWTOutput struct {
	AccessToken string `json:"access_token" valid:"-"`
}


type OTPInput struct {
	Token string `json:"token" valid:"length(6),numeric"`
}

type OTPOutput struct {
	Base32     string `json:"base32" valid:"required"`
	OtpauthUrl string `json:"otpauth_url" valid:"required"`
}
