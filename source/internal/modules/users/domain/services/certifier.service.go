package services

import (
	"certifications/internal/modules/users/domain/entities"
	"certifications/internal/modules/users/infra/repositories"
	"certifications/pkg/cache"

	errs "github.com/Lucasvmarangoni/logella/err"
)

type CertifierService struct {
	CertifierRepository repositories.CertifierRepository
	Memcached           cache.Mencacher[*entities.Certifier]
}

func NewCertifierService(
	certifierRepository repositories.CertifierRepository,
	memcached cache.Mencacher[*entities.Certifier],
) *CertifierService {
	return &CertifierService{
		CertifierRepository: certifierRepository,
		Memcached:           memcached,
	}
}

func (c *CertifierService) setToMemcacheIfNotNil(certifier *entities.Certifier) error {
	err := c.Memcached.Set(certifier.User_ID.String(), certifier)
	if err != nil {
		errs.Assertion(err).Context("c.Memcached.Set")
	}
	return nil
}
