package database

import (
	"context"

	errs "github.com/Lucasvmarangoni/logella/err"
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
		return errs.Wrap(err, "r.initUserTable", errs.GetHTTPStatusFromPgError(err))
	}

	err = r.initCertifierTable(ctx)
	if err != nil {
		return errs.Wrap(err, "r.initCertifierTable", errs.GetHTTPStatusFromPgError(err))
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
			otp_secret BYTEA,
			otp_auth_url BYTEA,			
			otp_enabled BOOL,
			created_at TIMESTAMP NOT NULL			
		)`)
	if err != nil {
		return errs.Wrap(err, "r.tx.Exec", errs.GetHTTPStatusFromPgError(err))
	}
	log.Info().Str("context", "TableRepository").Msg("Database - Created users table successfully.")
	return nil
}

func (r *TableRepositoryDb) initCertifierTable(ctx context.Context) error {

	_, err := r.tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS certifiers (
			user_id UUID PRIMARY KEY NOT NULL,
			description VARCHAR(5000) NOT NULL,
			validated BOOL NOT NULL DEFAULT FALSE,
			balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
			reserved NUMERIC(15,2) NOT NULL DEFAULT 0.00,
			account_status SMALLINT NOT NULL	
		)`)
	if err != nil {
		return errs.Wrap(err, "r.tx.Exec", errs.GetHTTPStatusFromPgError(err))
	}
	log.Info().Str("context", "TableRepository").Msg("Database - Created certifers table successfully.")
	return nil
}
