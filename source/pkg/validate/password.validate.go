package validate

import (
	"errors"
	"regexp"
)

// Using the golang standard error type because it will be sent in the response
func ValidatePassword(password string) error {
	if len(password) < 10 {
		return errors.New("password must be at least 10 characters long")
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one digit")
	}

	if !regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
