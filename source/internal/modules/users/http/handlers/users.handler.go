package handlers

import (
	"certifications/internal/modules/users/domain/services"
	"context"
	go_err "errors"
	"regexp"
)

type UserHandler struct {
	userService *services.UserService
	ctx         context.Context
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) validateCPF(cpf *string) error {
	cpfRegex := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	if *cpf != "" && !cpfRegex.MatchString(*cpf) {
		return go_err.New("invalid CPF format")
	}
	return nil
}

func (u *UserHandler) validateEmail(email *string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if *email != "" && !emailRegex.MatchString(*email) {
		return go_err.New("invalid email format")
	}
	return nil
}

func (u *UserHandler) validateNameAndLastname(name, lastname *string) error {
	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`)

	if *name != "" && !nameRegex.MatchString(*name) {
		return go_err.New("invalid name format")
	}
	if *lastname != "" && !nameRegex.MatchString(*lastname) {
		return go_err.New("invalid last name format")
	}
	return nil
}
