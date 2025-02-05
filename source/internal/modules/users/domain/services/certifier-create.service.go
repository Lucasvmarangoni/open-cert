package services

import (
	"certifications/internal/modules/users/domain/entities"
	pkg_entities "certifications/pkg/entities"
	"context"

	errs "github.com/Lucasvmarangoni/logella/err"
)

func (c *CertifierService) CreateCertifier(id, description string) error {
	user_id, err := pkg_entities.ParseID(id)
	if err != nil {
		return errs.Assertion(err).Context("pkg_entities.ParseID")
	}
	certifier, err := entities.NewCertifier(user_id, description)
	if err != nil {
		return errs.Assertion(err).Context("entities.NewCertifier")
	}
	err = c.CertifierRepository.Insert(certifier, context.Background())
	if err != nil {
		return errs.Assertion(err).Context("c.CertifierRepository.Insert")
	}
	err = c.setToMemcacheIfNotNil(certifier)
	if err != nil {
		return errs.Assertion(err).Context("c.setToMemcacheIfNotNil")
	}
	return nil
}
