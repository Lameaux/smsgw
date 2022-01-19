package inputs

import "time"

type SearchParams struct {
	Offset int
	Limit  int

	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}
