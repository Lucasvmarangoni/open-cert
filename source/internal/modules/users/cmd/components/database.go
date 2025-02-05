package components

import (
	database_connection "certifications/internal/infra"
	"certifications/internal/modules/users/infra/database"
	"context"
	"net/http"

	errs "github.com/Lucasvmarangoni/logella/err"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

func Database(ctx context.Context, db database_connection.Config) (*pgx.Conn, error) {
	conn, err := db.Connect(ctx)
	if err != nil {
		return nil, errs.Wrap(err, "db.Connect", http.StatusInternalServerError)
	}

	err = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		repo := database.NewTableRepository(tx)
		return repo.InitTables(ctx)
	})
	if err != nil {
		return nil, errs.Wrap(err, "repo.InitTables", http.StatusInternalServerError)
	}
	return conn, nil
}
