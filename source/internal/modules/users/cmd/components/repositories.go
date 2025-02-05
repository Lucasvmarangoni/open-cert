package components

import (
	"certifications/config"
	"certifications/internal/modules/users/infra/repositories"

	"github.com/jackc/pgx/v5"
)

type UsersRespositoryFactory struct {
	UserRepository      *repositories.UserRepositoryDb
	CertifierRepository *repositories.CertifierRepositoryDb
}

func NewRepositoriesFactory(conn *pgx.Conn) *UsersRespositoryFactory {
	userRepository := &UsersRespositoryFactory{
		UserRepository:      repositories.NewUserRepository(conn, config.GetEnvString("database", "sym_encrypt")),
		CertifierRepository: repositories.NewCertifierRepository(conn),
	}
	return userRepository
}
