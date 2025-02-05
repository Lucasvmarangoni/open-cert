package entities

import (
	"certifications/pkg/entities"
	"certifications/pkg/validate"
	"net/http"
	"time"

	errs "github.com/Lucasvmarangoni/logella/err"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            entities.ID         `json:"id" valid:"required"`
	Name          string              `json:"name" valid:"length(3|10),matches(^[a-zA-Z ]+$)"`
	LastName      string              `json:"last_name" valid:"length(3|50),matches(^[a-zA-Z ]+$)"`
	CPF           string              `json:"cpf" valid:"matches(^[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}-[0-9]{2}$)"`
	Email         string              `json:"email" valid:"required,email"`
	Password      string              `json:"password" valid:"required"`
	AccountStatus entities.UserStatus `json:"account_status" valid:"required"`
	Certifier     bool                `json:"certifier" valid:"-"`
	OtpSecret     string              `json:"otp_secret" valid:"-"`
	OtpAuthUrl    string              `json:"otp_auth_url" valid:"-"`
	OtpEnabled    bool                `json:"otp_enabled" valid:"-"`
	CreatedAt     time.Time           `json:"created_at" valid:"required"`
}

var user_status = entities.User_GetStatus()

func (u *User) ValidateHashPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errs.Wrap(err, "bcrypt.CompareHashAndPassword", http.StatusBadRequest)
	}
	return nil
}

func NewUser(name, lastName, cpf, email, password string) (*User, error) {

	hashPassword, err := bcryptHash(password)
	if err != nil {
		return nil, errs.Wrap(err, "u.PreparePassword", http.StatusBadRequest)
	}

	user := &User{
		Name:          name,
		LastName:      lastName,
		CPF:           cpf,
		Email:         email,
		Password:      hashPassword,
		AccountStatus: entities.UserStatus(user_status.GetRegular()),
		Certifier:     false,
	}
	user.prepare()

	err = Validate(user)
	if err != nil {
		return nil, errs.Wrap(err, "user.Validate", http.StatusBadRequest)
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

func (u *User) Update(email, password string) error {
	u.Email = email

	hashPassword, err := bcryptHash(password)
	if err != nil {
		return errs.Wrap(err, "u.PreparePassword", http.StatusBadRequest)
	}
	u.Password = string(hashPassword)

	err = Validate(u)
	if err != nil {
		return errs.Wrap(err, "Reputation.Validate", http.StatusBadRequest)
	}
	return nil
}

func bcryptHash(input string) (string, error) {
	err := validate.ValidatePassword(input)
	if err != nil {
		return "", errs.Wrap(err, "validatePassword", http.StatusBadRequest)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.Wrap(err, "bcrypt.GenerateFromPassword", http.StatusBadRequest)
	}
	return string(hashPassword), nil
}
