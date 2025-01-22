package entities

import (
	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/google/uuid"
)

type ID = uuid.UUID

func NewID() ID {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		errs.Ctx(err, "uuid.NewV7()")
	}
	return uuidV7
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

func Nil() uuid.UUID {
	return uuid.Nil
}
