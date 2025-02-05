package database

import (
	"context"
	"os"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Db       *pgx.ConnConfig
	DbName   string
	User     string
	Password string
	SSLMode  string
	Port     string
}

func NewDb() *Config {
	return &Config{}
}

func NewDbTest(ctx context.Context) *pgx.Conn {
	dbInstance := NewDb()
	dbInstance.DbName = "file-manager"
	dbInstance.Port = os.Getenv("DATABASE_PORT")
	dbInstance.User = os.Getenv("DATABASE_USER")
	dbInstance.Password = os.Getenv("DATABASE_PASSWORD")
	dbInstance.SSLMode = os.Getenv("DATABASE_SSL_MODE")

	connection, err := dbInstance.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err).Str("file", "db.go").Str("Method", "NewDbTest").Str("Operation", "dbInstance.Connect").Msg("Error in the configuration of the database testing")
	}
	return connection
}

func (cfg *Config) Connect(ctx context.Context) (*pgx.Conn, error) {
	var err error
	url := "postgresql://" + cfg.User + /* ":" + cfg.Password */ "@roach1:26257,roach2:26258,roach3:26259/" + cfg.DbName + "?sslmode=" + cfg.SSLMode

	cfg.Db, err = pgx.ParseConfig(url)
	if err != nil {
		return nil, errs.Wrap(err, "pgx.ParseConfig(url)", errs.GetHTTPStatusFromPgError(err))
	}

	cfg.Db.RuntimeParams["application_name"] = "financial_file_manager"
	conn, err := pgx.ConnectConfig(ctx, cfg.Db)
	if err != nil {
		return nil, errs.Wrap(err, "pgx.ConnectConfig", errs.GetHTTPStatusFromPgError(err))
	}
	return conn, nil
}
