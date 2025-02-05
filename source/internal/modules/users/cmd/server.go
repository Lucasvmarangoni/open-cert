package main

import (
	"certifications/config"
	database_connection "certifications/internal/infra"
	"certifications/internal/modules/users/cmd/components"
	"context"
	"net/http"
	"time"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/go-chi/chi"

	"github.com/Lucasvmarangoni/logella/config/log"
)

var db database_connection.Config

func init() {
	logger.ConfigDefault()

	db.DbName = config.GetEnvString("database", "name")
	db.Port = config.GetEnvString("database", "port")
	db.User = config.GetEnvString("database", "user")
	db.Password = config.ReadSecretString(config.GetEnvString("database", "password"))
	db.SSLMode = config.GetEnvString("database", "ssl_mode")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := components.Database(ctx, db)
	errs.FailOnErrLog(err, "Database(ctx)", "Failed exec Database")

	repository := components.NewRepositoriesFactory(conn)

	queueFactory := components.NewQueueFactory(ctx, repository.UserRepository)
	defer queueFactory.Ch.Close()

	cacheFactory := components.NewCacheFactory()

	services := components.NewServiceFactory(
		repository,
		cacheFactory,
		queueFactory.RabbitMQ,
		queueFactory.CreateMessageChannel,
		queueFactory.CreateReturnChannel,
	)
	r := chi.NewRouter()

	components.Http(r, services.UserService, services.CertifierService)

	err = http.ListenAndServe(":8000", r)
	errs.FailOnErrLog(err, "http.ListenAndServe", "Failed server listen")
}
