package entities

import (
	"certifications/pkg/entities"

	"github.com/Lucasvmarangoni/logella/err"
)

type Reputation struct {
	ID      entities.ID `json:"id" valid:"required"`
	User_ID entities.ID `json:"user_id" valid:"required"`
	Rating  int         `json:"rating" valid:"required,range(1|5)"`
	Comment string      `json:"comment" valid:"required,length(0|250)"`
}

func NewReputation(user_id entities.ID, rating int, comment string) (*Reputation, error) {

	reputation := &Reputation{
		ID:      entities.NewID(),
		User_ID: user_id,
		Rating:  rating,
		Comment: comment,
	}

	err := Validate(reputation)
	if err != nil {
		return nil, errs.Ctx(err, "reputation.Validate")
	}
	return reputation, nil
}
