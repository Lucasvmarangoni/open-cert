package repositories

import (
	"certifications/internal/modules/users/domain/entities"
	pkg_entities "certifications/pkg/entities"
	"context"

	errs "github.com/Lucasvmarangoni/logella/err"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

type ReputationRepository interface {
	Get_AVG(user_id pkg_entities.ID, ctx context.Context) (float64, error)
	Insert(reputatation entities.Reputation, ctx context.Context) error
	Update(id pkg_entities.ID, rating int, comment string, ctx context.Context) error
	Delete(id string, ctx context.Context) error
}

type ReputationRepositoryDb struct {
	conn *pgx.Conn
}

func NewReputationRepository(conn *pgx.Conn) *ReputationRepositoryDb {
	return &ReputationRepositoryDb{
		conn: conn,
	}
}

func (r *ReputationRepositoryDb) Insert(reputatation entities.Reputation, ctx context.Context) error {
	sql := `INSERT INTO repuations (id, user_id, rating, comment) VALUES ($1, $2, $3, $4)`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {

		_, err := tx.Exec(ctx, sql,
			reputatation.ID,
			reputatation.User_ID,
			reputatation.Rating,
			reputatation.Comment,
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

func (r *ReputationRepositoryDb) Get_AVG(user_id pkg_entities.ID, ctx context.Context) (float64, error) {
	sql := `SELECT AVG(rating) AS average_rating 
	FROM reputations 
	WHERE user_id = $1`

	var averageRating float64

	err := r.conn.QueryRow(ctx, sql, user_id).Scan(&averageRating)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, errs.Wrap(err, "QueryRow.Scan", errs.GetHTTPStatusFromPgError(err))
	}
	return averageRating, nil
}

func (r *ReputationRepositoryDb) Update(id pkg_entities.ID, rating int, comment string, ctx context.Context) error {
	sql := `UPDATE reputations SET rating = $2, comment = $3, WHERE id = $1`
	err := crdbpgx.ExecuteTx(ctx, r.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sql,
			id,
			rating,
			comment,
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

func (r *ReputationRepositoryDb) Delete(id string, ctx context.Context) error {
	sql := `DELETE FROM reputations WHERE id = $1`
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
