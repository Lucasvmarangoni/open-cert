package repositories

import (
	"certifications/internal/modules/users/domain/entities"
	pkg_entities "certifications/pkg/entities"
	"context"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Insert(user *entities.User, ctx context.Context) error
	FindByEmail(email string, ctx context.Context) (*entities.User, error)
	FindById(id pkg_entities.ID, ctx context.Context) (*entities.User, error)
	FindByCpf(cpf string, ctx context.Context) (*entities.User, error)
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
	VALUES ($1, $2, encrypt($3, $13, 'aes/pad:pkcs'), encrypt($4, $13, 'aes/pad:pkcs'), encrypt($5, $13, 'aes/pad:pkcs'), $6, $7, $8, $9, $10, $11, $12)`

	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			user.ID,
			user.Name,
			[]byte(user.LastName),
			[]byte(user.CPF),
			[]byte(user.Email),
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
			return errs.Ctx(err, "r.tx.Exec")
		}
		return nil
	})
	if err != nil {
		return errs.Ctx(err, "crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) FindById(id pkg_entities.ID, ctx context.Context) (*entities.User, error) {
	sql := `SELECT * FROM users WHERE id = $1`
	var row pgx.Row
	user := &entities.User{}
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		row = tx.QueryRow(ctx, sql, id)
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
			return errs.Ctx(err, "row.Scan")
		}
		return nil
	})
	if err != nil {
		return nil, errs.Ctx(err, "crdbpgx.ExecuteTx")
	}
	return user, nil
}

func (r *UserRepositoryDb) FindByEmail(hashEmail string, ctx context.Context) (*entities.User, error) {
	sql := `SELECT * FROM users WHERE hash_email = $1`
	var row pgx.Row
	user := &entities.User{}
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		row = tx.QueryRow(ctx, sql, hashEmail)
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
			return errs.Ctx(err, "row.Scan")
		}
		return nil
	})
	if err != nil {
		return nil, errs.Ctx(err, "crdbpgx.ExecuteTx")
	}
	return user, nil
}

func (r *UserRepositoryDb) FindByCpf(hashCPF string, ctx context.Context) (*entities.User, error) {
	sql := `SELECT * FROM users WHERE hash_cpf = $1`
	var row pgx.Row
	user := &entities.User{}
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		row = tx.QueryRow(ctx, sql, hashCPF)
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
			return errs.Ctx(err, "row.Scan")
		}
		return nil
	})
	if err != nil {
		return nil, errs.Ctx(err, "crdbpgx.ExecuteTx")
	}
	return user, nil
}

func (r *UserRepositoryDb) Update(user *entities.User, ctx context.Context) error {
	sql := `UPDATE users SET email = $2, password = $3, update_log = $6 WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			user.ID,
			user.Email,
			user.Password,
		)
		if err != nil {
			return errs.Ctx(err, "tx.Exec(")
		}
		return nil
	})
	if err != nil {
		return errs.Ctx(err, "crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) UpdateOTP(user *entities.User, ctx context.Context) error {
	sql := `UPDATE users SET otp_auth_url = $2, otp_secret = $3 WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			user.ID,
			user.OtpAuthUrl,
			user.OtpSecret,
		)
		if err != nil {
			return errs.Ctx(err, "tx.Exec")
		}
		return nil
	})
	if err != nil {
		return errs.Ctx(err, "crdbpgx.ExecuteTx")
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
			return errs.Ctx(err, "tx.Exec")
		}
		return nil
	})
	if err != nil {
		return errs.Ctx(err, "crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) Delete(id string, ctx context.Context) error {
	sql := `DELETE FROM users WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql, id)
		if err != nil {
			return errs.Ctx(err, "tx.Exec(")
		}
		return nil
	})
	if err != nil {
		return errs.Ctx(err, "crdbpgx.ExecuteTx")
	}
	return nil
}

func (r *UserRepositoryDb) CheckIfUserAlreadyExists(hashEmail, hashCPF string, ctx context.Context) (bool, error) {
	sql := `SELECT COUNT(*) FROM users WHERE hash_email = $1 OR hash_cpf = $2`
	var count int
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, sql, hashEmail, hashCPF)
		err := row.Scan(&count)
		if err != nil {
			return errs.Ctx(err, "row.Scan")
		}
		return nil
	})
	if err != nil {
		return false, errs.Ctx(err, "crdbpgx.ExecuteTx")
	}
	return count > 0, nil
}
