package repositories

import (
	"certifications/internal/modules/users/domain/entities"
	pkg_entities "certifications/pkg/entities"
	"context"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Insert(user *entities.User, ctx context.Context) error
	FindByEmail(email string, ctx context.Context) (*entities.User, error)
	FindById(id pkg_entities.ID, ctx context.Context) (*entities.User, error)
	Update(user *entities.User, ctx context.Context) error
	UpdateOTP(user *entities.User, ctx context.Context) error
	UpdateOTPVerify(user *entities.User, ctx context.Context) error
	Delete(id string, ctx context.Context) error
	CheckIfUserAlreadyExists(email, cpf string, ctx context.Context) (bool, error)
}

type UserRepositoryDb struct {
	conn *pgx.Conn
	key  string
}

func NewUserRepository(conn *pgx.Conn, key string) *UserRepositoryDb {
	return &UserRepositoryDb{
		conn: conn,
		key:  key,
	}
}

func (r *UserRepositoryDb) Insert(user *entities.User, ctx context.Context) error {
	if user.ID.String() == "" {
		user.ID = pkg_entities.NewID()
	}

	sql := `INSERT INTO users (id, name, last_name, cpf, email, password, account_status, certifier, otp_secret, otp_auth_url, otp_enabled, created_at) 
	VALUES ($1, $2, encrypt($3::BYTES, $13::BYTES, 'aes'), encrypt($4::BYTES, $13::BYTES, 'aes'), 
	encrypt($5::BYTES, $13::BYTES, 'aes'), $6, $7, $8, encrypt($9::BYTES, $13::BYTES, 'aes'), encrypt($10::BYTES, $13::BYTES, 'aes'), $11, $12)`

	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			user.ID,
			user.Name,
			user.LastName,
			user.CPF,
			user.Email,
			user.Password,
			user.AccountStatus,
			user.Certifier,
			user.OtpSecret,
			user.OtpAuthUrl,
			user.OtpEnabled,
			user.CreatedAt,
			r.key,
		)
		if err != nil {
			return errs.Wrap(err, "r.tx.Exec", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) FindById(id pkg_entities.ID, ctx context.Context) (*entities.User, error) {
	sql := `SELECT 
	id,
    name,
	convert_from(decrypt(last_name, $2::BYTES, 'aes'), 'UTF8') AS last_name,
	convert_from(decrypt(cpf, $2::BYTES, 'aes'), 'UTF8') AS cpf,
	convert_from(decrypt(email, $2::BYTES, 'aes'), 'UTF8') AS email,
	password,
    account_status,
    certifier,
	convert_from(decrypt(otp_secret::BYTES, $2::BYTES, 'aes'), 'UTF8') AS otp_secret, 
	convert_from(decrypt(otp_auth_url::BYTES, $2::BYTES, 'aes'), 'UTF8') AS otp_auth_url,
	otp_enabled,
    created_at
	FROM users 
	WHERE id = $1;
    `
	user := &entities.User{}
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, sql, id, r.key)
		err := row.Scan(
			&user.ID,
			&user.Name,
			&user.LastName,
			&user.CPF,
			&user.Email,
			&user.Password,
			&user.AccountStatus,
			&user.Certifier,
			&user.OtpSecret,
			&user.OtpAuthUrl,
			&user.OtpEnabled,
			&user.CreatedAt,
		)
		if err != nil {
			return errs.Wrap(err, "row.Scan", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return nil, errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return user, nil
}

func (r *UserRepositoryDb) FindByEmail(email string, ctx context.Context) (*entities.User, error) {

	sql := `SELECT 
	id,
    name,
	convert_from(decrypt(last_name, $2::BYTES, 'aes'), 'UTF8') AS last_name,
	convert_from(decrypt(cpf, $2::BYTES, 'aes'), 'UTF8') AS cpf,
	convert_from(decrypt(email, $2::BYTES, 'aes'), 'UTF8') AS email,
	password,
    account_status,
    certifier,
	convert_from(decrypt(otp_secret::BYTES, $2::BYTES, 'aes'), 'UTF8') AS otp_secret, 
	convert_from(decrypt(otp_auth_url::BYTES, $2::BYTES, 'aes'), 'UTF8') AS otp_auth_url,
	otp_enabled,
    created_at
	FROM users 
	WHERE convert_from(decrypt(email, $2::BYTES, 'aes'), 'UTF8') = $1;
    `
	var row pgx.Row
	user := &entities.User{}
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		row = tx.QueryRow(ctx, sql, email, r.key)
		err := row.Scan(
			&user.ID,
			&user.Name,
			&user.LastName,
			&user.CPF,
			&user.Email,
			&user.Password,
			&user.AccountStatus,
			&user.Certifier,
			&user.OtpSecret,
			&user.OtpAuthUrl,
			&user.OtpEnabled,
			&user.CreatedAt,
		)
		if err != nil {
			return errs.Wrap(err, "row.Scan", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return nil, errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return user, nil
}

func (r *UserRepositoryDb) Update(user *entities.User, ctx context.Context) error {
	sql := `UPDATE users SET email = encrypt($2::BYTES, $4::BYTES, 'aes'), password = $3 WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			user.ID,
			user.Email,
			user.Password,
			r.key,
		)
		if err != nil {
			return errs.Wrap(err, "tx.Exec", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) UpdateOTP(user *entities.User, ctx context.Context) error {
	sql := `UPDATE users SET otp_auth_url = encrypt($2::BYTES, $4::BYTES, 'aes'), otp_secret = encrypt($3::BYTES, $4::BYTES, 'aes') WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			user.ID,
			user.OtpAuthUrl,
			user.OtpSecret,
			r.key,
		)
		if err != nil {
			return errs.Wrap(err, "tx.Exec", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) UpdateOTPVerify(user *entities.User, ctx context.Context) error {
	sql := `UPDATE users SET otp_enabled = $2 WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			user.ID,
			user.OtpEnabled,
		)
		if err != nil {
			return errs.Wrap(err, "tx.Exec", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) Delete(id string, ctx context.Context) error {
	sql := `DELETE FROM users WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, id)
		if err != nil {
			return errs.Wrap(err, "tx.Exec", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) CheckIfUserAlreadyExists(email, cpf string, ctx context.Context) (bool, error) {
	sql := `SELECT COUNT(*) FROM users WHERE decrypt(email, $3, 'aes') = $1 OR decrypt(cpf, $3, 'aes') = $2`
	var count int
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, sql, email, cpf, r.key)
		err := row.Scan(&count)
		if err != nil {
			return errs.Wrap(err, "row.Scan", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return false, errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return count > 0, nil
}
