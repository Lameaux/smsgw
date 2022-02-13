package repos

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

const (
	connTimeout  = 1 * time.Second
	queryTimeout = 3 * time.Second
	txTimeout    = 1 * time.Second
)

func DBConnContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), connTimeout)
}

func DBQueryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), queryTimeout)
}

func DBTxContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), txTimeout)
}

func DBQueryGet(conn db.Conn, dst interface{}, query string, args ...interface{}) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	err := pgxscan.Get(ctx, conn, dst, query, args...)
	if pgxscan.NotFound(err) {
		return models.ErrNotFound
	}

	return err
}

func DBQuerySelect(conn db.Conn, dst interface{}, query string, args ...interface{}) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	return pgxscan.Select(ctx, conn, dst, query, args...)
}

func DBQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func appendMessageParams(q *inputs.MessageParams, stmt string, args []interface{}) (string, []interface{}) {
	if q.MSISDN != nil {
		args = append(args, q.MSISDN)
		stmt += fmt.Sprintf("and msisdn = $%d\n", len(args))
	}

	if q.Status != nil {
		args = append(args, q.Status)
		stmt += fmt.Sprintf("and status = $%d\n", len(args))
	}

	return stmt, args
}

func appendSearchParams(q *inputs.SearchParams, sb sq.SelectBuilder) sq.SelectBuilder {
	if q.CreatedAtFrom != nil {
		sb = sb.Where("and created_at >= ?", q.CreatedAtFrom)
	}

	if q.CreatedAtTo != nil {
		sb = sb.Where("and created_at <= ?", q.CreatedAtTo)
	}

	return sb.OrderBy("created_at desc").Limit(q.Limit).Offset(q.Offset)
}
