package entities

import (
	"fmt"
	"net/http"
	"time"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/asaskevich/govalidator"
)

func isStartBeforeEnding(start, ending time.Time) error {
	before := start.Before(ending)

	if before {
		return nil
	}
	return fmt.Errorf("Start is not earlier than Ending")
}

func Validate[T any](data T) error {
	_, err := govalidator.ValidateStruct(data)
	if err != nil {
		return errs.Wrap(err, "govalidator.ValidateStruct", http.StatusBadRequest)
	}
	return nil
}
