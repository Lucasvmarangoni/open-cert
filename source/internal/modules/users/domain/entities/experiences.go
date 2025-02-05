package entities

import (
	"certifications/pkg/entities"
	"net/http"
	"time"

	errs "github.com/Lucasvmarangoni/logella/err"
)

type Experiences struct {
	ID          entities.ID `json:"id" valid:"required"`
	UserID      entities.ID `json:"user_id" valid:"required"`
	Title       string      `json:"title" valid:"required,length(3|50)"`
	Company     string      `json:"company" valid:"required,length(3|50)"`
	Description string      `json:"description" valid:"required,length(3|2500)"`
	Start       time.Time   `json:"start" valid:"required"`
	Ending      time.Time   `json:"ending" valid:"required"`
}

func (e *Experiences) NewExperience(user_id entities.ID, title, company, description string, start, ending time.Time) (*Experiences, error) {
	err := isStartBeforeEnding(start, ending)
	if err != nil {
		return nil, errs.Wrap(err, "isStartBeforeEnding", http.StatusBadRequest)
	}

	experience := Experiences{
		ID:          entities.NewID(),
		UserID:      user_id,
		Title:       title,
		Company:     company,
		Description: description,
		Start:       start,
		Ending:      ending,
	}

	err = Validate(experience)
	if err != nil {
		return nil, errs.Assertion(err).Context("experience.Validate")
	}
	return &experience, nil
}
