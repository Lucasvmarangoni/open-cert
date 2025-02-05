package repositories

import (
	"certifications/internal/modules/users/domain/entities"
	pkg_entities "certifications/pkg/entities"
	"context"

	errs "github.com/Lucasvmarangoni/logella/err"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

type CertifierRepository interface {
	Insert(certifier *entities.Certifier, ctx context.Context) error
	Find_ById(id pkg_entities.ID, ctx context.Context) (*entities.Certifier, error)
	Update_Description(certifier entities.Certifier, ctx context.Context) error
	Delete(id string, ctx context.Context) error
}

type CertifierRepositoryDb struct {
	conn *pgx.Conn
}

func NewCertifierRepository(conn *pgx.Conn) *CertifierRepositoryDb {
	return &CertifierRepositoryDb{
		conn: conn,
	}
}

func (c *CertifierRepositoryDb) Insert(certifier *entities.Certifier, ctx context.Context) error {

	sql := `INSERT INTO certifiers (user_id, description, validated, balance, reserved, account_status) VALUES ($1, $2, $3, $4, $5, $6)`
	err := crdbpgx.ExecuteTx(ctx, c.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {

		_, err := tx.Exec(ctx, sql,
			certifier.User_ID,
			certifier.Description,
			certifier.Validated,
			certifier.Balance,
			certifier.Reserved,
			certifier.AccountStatus,
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

func (c *CertifierRepositoryDb) Find_ById(id pkg_entities.ID, ctx context.Context) (*entities.Certifier, error) {
	sql := `SELECT * FROM certifiers WHERE user_id = $1`
	var row pgx.Row
	certifier := &entities.Certifier{}
	err := crdbpgx.ExecuteTx(ctx, c.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		row = tx.QueryRow(ctx, sql, id)
		err := row.Scan(
			&certifier.User_ID,
			&certifier.Description,
			&certifier.Validated,
			&certifier.Balance,
			&certifier.Reserved,
		)
		if err != nil {
			return errs.Wrap(err, "row.Scan", errs.GetHTTPStatusFromPgError(err))
		}
		return nil
	})
	if err != nil {
		return nil, errs.Assertion(err).Context("crdbpgx.ExecuteTx")
	}
	return certifier, nil
}

func (c *CertifierRepositoryDb) Update_Description(certifier entities.Certifier, ctx context.Context) error {
	sql := `UPDATE certifiers SET description = $2, 6 WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, c.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			certifier.User_ID,
			certifier.Description,
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

func (c *CertifierRepositoryDb) Delete(id string, ctx context.Context) error {
	// Criar uma tabela que registra as solicitações de exclusão dos dados de certificador.
	return nil
}

// COLOCAR NO MODULO DE CERTIFICATIONS
// func (c *CertifierRepositoryDb) Get_ApprovalRating(user_id pkg_entities.ID, ctx context.Context) (float64, error) {
// 	sql := `SELECT AVG(approved::int) AS approval_rating
// 	FROM certifications
// 	WHERE user_id = $1;`

// 	var approval_rating float64

// 	err := c.conn.QueryRow(ctx, sql, user_id).Scan(&approval_rating)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return 0, nil
// 		}
// 		return 0, errs.Wrap(err, "QueryRow.Scan")
// 	}
// 	return approval_rating, nil
// }
