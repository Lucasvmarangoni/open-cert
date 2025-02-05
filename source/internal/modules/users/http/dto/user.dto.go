package dto

import (
	"certifications/pkg/entities"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type UserInput struct {
	Name     string `json:"name" valid:"required,matches(^[a-zA-Z ]+$),length(3|10)"`
	LastName string `json:"last_name" valid:"required,matches(^[a-zA-Z ]+$),length(3|50)"`
	CPF      string `json:"cpf" valid:"matches(^[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}-[0-9]{2}$)"`
	Email    string `json:"email" valid:"email"`
	Password string `valid:"required,minstringlength(10),matches([a-z]),matches([A-Z]),matches([0-9]),matches([!@#\\$%\\^&\\*])"`
}

type AuthenticationInput struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `valid:"required,minstringlength(10),matches([a-z]),matches([A-Z]),matches([0-9]),matches([!@#\\$%\\^&\\*])"`
}

type GetJWTOutput struct {
	AccessToken      string `json:"access_token" valid:"required"`
	TwoAuthnRequirid bool   `json:"two_authn_required" valid:"-"`
}

type OTPInput struct {
	Token string `json:"token" valid:"length(6),numeric"`
}

type OTPOutput struct {
	Base32      string `json:"base32" valid:"required"`
	OtpauthUrl  string `json:"otpauth_url" valid:"required"`
	AccessToken string `json:"access_token" valid:"required"`
}

type UserOutput struct {
	ID        entities.ID `json:"id" valid:"-"`
	Name      string      `json:"name" valid:"required,matches(^[a-zA-Z ]+$),length(3|10)"`
	LastName  string      `json:"last_name" valid:"required,matches(^[a-zA-Z ]+$),length(3|50)"`
	CPF       string      `json:"cpf" valid:"required,matches(^[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}-[0-9]{2}$)"`
	Email     string      `json:"email" valid:"email"`
	Certifier bool        `json:"certifier" valid:"required"`
}

type UserUpdateInput struct {
	Password    string `valid:"required,minstringlength(10),matches([a-z]),matches([A-Z]),matches([0-9]),matches([!@#\\$%\\^&\\*])"`
	NewEmail    string `json:"email" valid:"email"`
	NewPassword string `json:"new_password" valid:"-"`
}

type CertifierInput struct {
	Description string `json:"description" valid:"length(200|5000),required"`
}
