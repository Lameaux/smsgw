package repos

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

const (
	connTimeout  = 1 * time.Second
	queryTimeout = 3 * time.Second
	txTimeout    = 1 * time.Second
)

type sqler interface {
	ToSql() (string, []interface{}, error)
}

func DBConnContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), connTimeout)
}

func DBTxContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), txTimeout)
}

func dbQueryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), queryTimeout)
}

func Begin(conn db.Conn) (pgx.Tx, error) { //nolint:ireturn
	ctx, cancel := DBTxContext()
	defer cancel()

	return conn.Begin(ctx)
}

func Rollback(tx pgx.Tx) error {
	ctx, cancel := DBTxContext()
	defer cancel()

	return tx.Rollback(ctx)
}

func Commit(tx pgx.Tx) error {
	ctx, cancel := DBTxContext()
	defer cancel()

	return tx.Commit(ctx)
}

func dbQuerySingle(conn db.Conn, dst interface{}, sb sqler) error {
	stmt, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	ctx, cancel := dbQueryContext()
	defer cancel()

	err = pgxscan.Get(ctx, conn, dst, stmt, args...)
	if pgxscan.NotFound(err) {
		return models.ErrNotFound
	}

	return err
}

func dbQueryAll(conn db.Conn, dst interface{}, sb sqler) error {
	stmt, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	ctx, cancel := dbQueryContext()
	defer cancel()

	return pgxscan.Select(ctx, conn, dst, stmt, args...)
}

func dbExec(conn db.Conn, sb sqler) error {
	stmt, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	ctx, cancel := dbQueryContext()
	defer cancel()

	_, err = conn.Exec(ctx, stmt, args...)

	return err
}

func dbQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func appendMessageParams(q *inputs.MessageParams, sb sq.SelectBuilder) sq.SelectBuilder {
	if q.MSISDN != nil {
		sb = sb.Where("msisdn = ?", q.MSISDN)
	}

	if q.Status != nil {
		sb = sb.Where("status = ?", q.Status)
	}

	return sb
}

func appendSearchParams(q *inputs.SearchParams, sb sq.SelectBuilder) sq.SelectBuilder {
	if q.CreatedAtFrom != nil {
		sb = sb.Where("created_at >= ?", q.CreatedAtFrom)
	}

	if q.CreatedAtTo != nil {
		sb = sb.Where("created_at <= ?", q.CreatedAtTo)
	}

	return sb.OrderBy("created_at desc").Limit(q.Limit).Offset(q.Offset)
}
