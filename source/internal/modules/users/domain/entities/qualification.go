package entities

import (
	"certifications/pkg/entities"
	"time"

	"github.com/Lucasvmarangoni/logella/err"
)

type Qualification struct {
	ID           entities.ID `json:"id" valid:"required"`
	UserID       entities.ID `json:"user_id" valid:"required"`
	Organization string      `json:"organization" valid:"required,length(3|50)"`
	Degree       string      `json:"degree" valid:"required,length(3|50)"`
	Start        time.Time   `json:"start" valid:"required"`
	Ending       time.Time   `json:"ending" valid:"required"`
}

func (q *Qualification) NewQualification(user_id entities.ID, organization, degree string, start, ending time.Time) (*Qualification, error) {
	err := isStartBeforeEnding(start, ending)
	if err != nil {
		return nil, errs.Ctx(err, "isStartBeforeEnding")
	}

	qualification := Qualification{
		ID:           entities.NewID(),
		UserID:       user_id,
		Organization: organization,
		Degree:       degree,
		Start:        start,
		Ending:       ending,
	}

	err = Validate(qualification)
	if err != nil {
		return nil, errs.Ctx(err, "qualification.Validate")
	}
	return &qualification, nil
}
