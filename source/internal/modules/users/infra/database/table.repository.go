package database

import (
	"context"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type TableRepository interface {
	InitTables(ctx context.Context) error
	createdPgcrypto(ctx context.Context) error
	initUserTable(ctx context.Context) error
}

type TableRepositoryDb struct {
	tx pgx.Tx
}

func NewTableRepository(db pgx.Tx) *TableRepositoryDb {
	return &TableRepositoryDb{tx: db}
}

func (r *TableRepositoryDb) InitTables(ctx context.Context) error {
	err := r.initUserTable(ctx)
	if err != nil {
		return errs.Ctx(err, "initUserTable")
	}

	err = r.createdPgcrypto(ctx)
	if err != nil {
		return errs.Ctx(err, "initUserTable")
	}
	return nil
}

func (r *TableRepositoryDb) initUserTable(ctx context.Context) error {

	_, err := r.tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY NOT NULL,
			name VARCHAR(15) NOT NULL,
			last_name BYTEA NOT NULL,
			cpf BYTEA UNIQUE NOT NULL,					
			email BYTEA UNIQUE NOT NULL,				
			password TEXT NOT NULL,		
			account_status SMALLINT NOT NULL,
			certifier BOOL NOT NULL DEFAULT FALSE,
			otp_secret TEXT,
			otp_auth_url TEXT,			
			otp_enabled BOOL,
			created_at TIMESTAMP NOT NULL			
		)`)
	if err != nil {
		return errs.Ctx(err, "r.tx.Exec")
	}
	log.Info().Str("context", "TableRepository").Msg("Database - Created users table successfully.")
	return nil
}

func (r *TableRepositoryDb) createdPgcrypto(ctx context.Context) error {

	_, err := r.tx.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS pgcrypto;
`)
	if err != nil {
		return errs.Ctx(err, "r.tx.Exec")
	}
	log.Info().Str("context", "TableRepository").Msg("Database - Created pgcrypto extension.")
	return nil
}
