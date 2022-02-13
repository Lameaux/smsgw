package inputs

import "time"

type SearchParams struct {
	Offset uint64
	Limit  uint64

	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}
