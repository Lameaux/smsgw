package repos

import (
	"context"
	"fmt"
	"time"

	"euromoby.com/smsgw/internal/inputs"
)

func DBConnContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 1*time.Second)
}

func DBQueryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func DBTxContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 1*time.Second)
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

func appendSearchParams(q *inputs.SearchParams, stmt string, args []interface{}) (string, []interface{}) {
	if q.CreatedAtFrom != nil {
		args = append(args, q.CreatedAtFrom)
		stmt += fmt.Sprintf("and created_at >= $%d\n", len(args))
	}

	if q.CreatedAtTo != nil {
		args = append(args, q.CreatedAtTo)
		stmt += fmt.Sprintf("and created_at <= $%d\n", len(args))
	}

	args = append(args, q.Limit, q.Offset)
	stmt += fmt.Sprintf("order by created_at desc limit $%d offset $%d", len(args)-1, len(args))

	return stmt, args
}
