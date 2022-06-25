package repos

import (
	sq "github.com/Masterminds/squirrel"

	"euromoby.com/smsgw/internal/inputs"
)

func AppendMessageParams(q *inputs.MessageParams, sb sq.SelectBuilder) sq.SelectBuilder {
	if q.MSISDN != nil {
		sb = sb.Where("msisdn = ?", q.MSISDN)
	}

	if q.Status != nil {
		sb = sb.Where("status = ?", q.Status)
	}

	return sb
}

func AppendSearchParams(q *inputs.SearchParams, sb sq.SelectBuilder) sq.SelectBuilder {
	if q.CreatedAtFrom != nil {
		sb = sb.Where("created_at >= ?", q.CreatedAtFrom)
	}

	if q.CreatedAtTo != nil {
		sb = sb.Where("created_at <= ?", q.CreatedAtTo)
	}

	return sb.OrderBy("created_at desc").Limit(q.Limit).Offset(q.Offset)
}
