package entities

import (
	"certifications/pkg/entities"
	"time"

	errs "github.com/Lucasvmarangoni/logella/err"
)

type UserUpdateLog struct {
	UserID    entities.ID            `json:"user_id" valid:"required"`
	Email     string                 `json:"email" valid:"-"`
	Password  string                 `json:"password" valid:"-"`
	Timestamp time.Time              `json:"timestamp" valid:"required"`
	Others    map[string]interface{} `json:"others" valid:"-"`
}

func NewUserUpdateLog(user_id entities.ID,
	email, password string,
	oldValues map[string]interface{}) (*UserUpdateLog, error) {

	user_update_log := &UserUpdateLog{
		UserID:   user_id,
		Email:    email,
		Password: password,
		Others:   oldValues,
	}

	err := Validate(user_update_log)
	if err != nil {
		return nil, errs.Assertion(err).Context("user_update_log.Validate")
	}

	return user_update_log, nil
}
