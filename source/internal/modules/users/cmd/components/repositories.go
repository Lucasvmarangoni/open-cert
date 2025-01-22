package components

import (
	"certifications/config"
	"certifications/internal/modules/users/infra/repositories"

	"github.com/jackc/pgx/v5"
)

type UsersRespositoryFactory struct {
	UserRepository *repositories.UserRepositoryDb
}

func NewRepositoriesFactory(conn *pgx.Conn) *UsersRespositoryFactory {
	userRepository := &UsersRespositoryFactory{
		UserRepository: repositories.NewUserRepository(conn, config.GetEnvString("database", "sym_encrypt")),
	}
	return userRepository
}
