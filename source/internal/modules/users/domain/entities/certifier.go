package entities

import (
	"certifications/pkg/entities"

	errs "github.com/Lucasvmarangoni/logella/err"
)

type Certifier struct {
	User_ID       entities.ID              `json:"id" valid:"required"`
	Description   string                   `json:"description" valid:"length(200|3000),required"`
	Validated     bool                     `json:"validated" valid:"-"`
	Balance       float64                  `json:"balance" valid:"-"`
	Reserved      float64                  `json:"reserved" valid:"-"`
	AccountStatus entities.CertifierStatus `json:"account_status" valid:"required"`
}

var certifier_status = entities.Certifier_GetStatus()

func NewCertifier(
	user_id entities.ID, description string) (*Certifier, error) {
	certifier := &Certifier{
		User_ID:       user_id,
		Description:   description,
		AccountStatus: entities.CertifierStatus(certifier_status.GetAnalysis()),
	}

	err := Validate(certifier)
	if err != nil {
		return nil, errs.Assertion(err).Context("certifier.Validate")
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

func (c *Certifier) UpdateAccountStatus(account_status entities.CertifierStatus) {
	c.AccountStatus = account_status
}
