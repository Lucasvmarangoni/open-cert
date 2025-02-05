package handlers

import (
	"certifications/internal/modules/users/domain/services"
	"context"
)

type UsersHandler struct {
	userService      *services.UserService
	certifierService *services.CertifierService
	ctx              context.Context
}

func NewUserHandler(userService *services.UserService, certifierService *services.CertifierService) *UsersHandler {
	return &UsersHandler{
		userService:      userService,
		certifierService: certifierService,
	}
}
