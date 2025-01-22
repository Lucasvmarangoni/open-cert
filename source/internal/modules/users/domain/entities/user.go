package entities

import (
	"certifications/pkg/entities"
	"certifications/pkg/validate"
	"fmt"
	"time"

	"github.com/Lucasvmarangoni/logella/err"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            entities.ID `json:"id" valid:"required"`
	Name          string      `json:"name" valid:"length(3|10),matches(^[a-zA-Z ]+$)"`
	LastName      string      `json:"last_name" valid:"length(3|50),matches(^[a-zA-Z ]+$)"`
	CPF           string      `json:"cpf" valid:"matches(^[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}-[0-9]{2}$)"`
	Email         string      `json:"email" valid:"email"`
	Password      string      `json:"password" valid:"required"`
	AccountStatus int         `json:"account_status" valid:"required"`
	Certifier     bool        `json:"certifier" valid:"-"`
	OtpSecret     string      `json:"otp_secret" valid:"-"`
	OtpAuthUrl    string      `json:"otp_auth_url" valid:"-"`
	OtpEnabled    bool        `json:"otp_enabled" valid:"-"`
	CreatedAt     time.Time   `json:"created_at" valid:"required"`
}

var status = entities.User_GetStatus()

func (u *User) ValidateStatus() error {

	switch u.AccountStatus {
	case status.GetBlocked(), status.GetRegular(), status.GetSuspend():
		return nil
	default:
		return fmt.Errorf("status inválido: %d", u.AccountStatus)
	}
}

func (u *User) ValidateHashPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errs.Ctx(err, "bcrypt.CompareHashAndPassword")
	}
	return nil
}

func NewUser(name, lastName, cpf, email, password string) (*User, error) {

	hashPassword, err := bcryptHash(password)
	if err != nil {
		return nil, errs.Ctx(err, "u.PreparePassword")
	}

	user := &User{
		Name:          name,
		LastName:      lastName,
		CPF:           cpf,
		Email:         email,
		Password:      hashPassword,
		AccountStatus: status.GetRegular(),
		Certifier:     false,
	}
	user.prepare()

	err = Validate(user)
	if err != nil {
		return nil, errs.Ctx(err, "user.Validate")
	}

	err = user.ValidateStatus()
	if err != nil {
		return nil, errs.Ctx(err, "user.ValidateStatus")
	}

	return user, nil
}

func (u *User) prepare() {
	if u.CreatedAt.IsZero() {
		u.ID = entities.NewID()
		u.CreatedAt = time.Now()
	}
}

func (u *User) PrepateTOTP(otpSecret, otpAuthUrl string) {
	u.OtpSecret = otpSecret
	u.OtpAuthUrl = otpAuthUrl
}

func (u *User) Update(name, lastName, email, password string) error {
	u.Name = name
	u.LastName = lastName
	u.Email = email

	hashPassword, err := bcryptHash(password)
	if err != nil {
		return errs.Ctx(err, "u.PreparePassword")
	}
	u.Password = string(hashPassword)

	err = Validate(u)
	if err != nil {
		return errs.Ctx(err, "Reputation.Validate")
	}

	return nil
}

func bcryptHash(input string) (string, error) {
	err := validate.ValidatePassword(input)
	if err != nil {
		return "", errs.Ctx(err, "validatePassword")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.Ctx(err, "bcrypt.GenerateFromPassword")
	}
	return string(hashPassword), nil
}
