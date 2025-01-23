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
