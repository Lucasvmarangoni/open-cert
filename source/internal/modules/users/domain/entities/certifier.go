package entities

import (
	"certifications/pkg/entities"

	errs "github.com/Lucasvmarangoni/logella/err"
)

type Certifier struct {
	User_ID       entities.ID `json:"id" valid:"required"`
	Description   string      `json:"description" valid:"length(3|3000),required"`
	Validated     bool        `json:"validated" valid:"required"`
	Balance       float64     `json:"balance" valid:"required"`
	Reserved      float64     `json:"reserved" valid:"required"`
	AccountStatus int         `json:"account_status" valid:"required"`
}

func NewCertifier(
	user_id entities.ID, description string) (*Certifier, error) {
	certifier := &Certifier{
		User_ID:     user_id,
		Description: description,
	}

	err := Validate(certifier)
	if err != nil {
		return nil, errs.Ctx(err, "certifier.Validate")
	}

	return certifier, nil
}

func (c *Certifier) UpdateBalance(balance float64) {
	c.Balance = balance
}

func (c *Certifier) UpdateReserved(reserved float64) {
	c.Reserved = reserved
}

func (c *Certifier) UpdateValidated(validated bool) {
	c.Validated = validated
}

func (c *Certifier) UpdateAccountStatus(account_status int) {
	c.AccountStatus = account_status
}
