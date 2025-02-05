package entities

import (
	"net/http"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/google/uuid"
)

type ID = uuid.UUID

func NewID() ID {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		errs.Wrap(err, "uuid.NewV7()", http.StatusInternalServerError)
	}
	return uuidV7
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), errs.Wrap(err, "uuid.Parse", http.StatusBadRequest)
}

func Nil() uuid.UUID {
	return uuid.Nil
}
