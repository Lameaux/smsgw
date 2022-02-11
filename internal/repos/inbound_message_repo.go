package repos

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

const (
	selectInboundMessagesBase = `select
	id, shortcode, status, msisdn, body,
	provider_id, provider_message_id,
	next_attempt_at, attempt_counter,
	created_at, updated_at
	from inbound_messages
	`

	constraintNameProviderMessageID = "inbound_provider_message_id"
)

type InboundMessageRepo struct {
	db db.Conn
}

func NewInboundMessageRepo(db db.Conn) *InboundMessageRepo {
	return &InboundMessageRepo{db}
}

func (r *InboundMessageRepo) Save(im *models.InboundMessage) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := `insert into inbound_messages (
		shortcode, status, msisdn, body,
		provider_id, provider_message_id,
		next_attempt_at, attempt_counter,
		created_at, updated_at
	)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	returning id
	`

	var insertedID string

	err := r.db.QueryRow(ctx, stmt,
		im.Shortcode,
		im.Status,
		im.MSISDN,
		im.Body,
		im.ProviderID,
		im.ProviderMessageID,
		im.NextAttemptAt,
		im.AttemptCounter,
		im.CreatedAt,
		im.UpdatedAt,
	).Scan(&insertedID)
	if err != nil {
		return r.wrapError(err)
	}

	im.ID = insertedID

	return nil
}

func (r *InboundMessageRepo) Update(im *models.InboundMessage) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	im.UpdatedAt = models.TimeNow()

	stmt := `update inbound_messages
	set status = $1, updated_at = $2
	where id = $3
	`

	_, err := r.db.Exec(ctx, stmt,
		im.Status,
		im.UpdatedAt,
		im.ID,
	)

	return err
}

func (r *InboundMessageRepo) FindByShortcodeAndID(shortcode, id string) (*models.InboundMessage, error) {
	stmt := selectInboundMessagesBase + "where shortcode = $1 AND id = $2"

	return r.querySingle(stmt, shortcode, id)
}

func (r *InboundMessageRepo) FindByQuery(q *inputs.InboundMessageSearchParams) ([]*models.InboundMessage, error) {
	stmt := selectInboundMessagesBase
	args := make([]interface{}, 0)

	args = append(args, q.Shortcode)
	stmt += fmt.Sprintf("where shortcode = $%d\n", len(args))

	stmt, args = appendMessageParams(q.MessageParams, stmt, args)
	stmt, args = appendSearchParams(q.SearchParams, stmt, args)

	return r.query(stmt, args...)
}

func (r *InboundMessageRepo) query(stmt string, args ...interface{}) ([]*models.InboundMessage, error) {
	messages := []*models.InboundMessage{}

	ctx, cancel := DBQueryContext()
	defer cancel()

	rows, err := r.db.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var om models.InboundMessage
		err = r.scanRow(rows, &om)

		if err != nil {
			return nil, err
		}

		messages = append(messages, &om)
	}

	return messages, nil
}

func (r *InboundMessageRepo) querySingle(stmt string, args ...interface{}) (*models.InboundMessage, error) {
	ctx, cancel := DBQueryContext()
	defer cancel()

	row := r.db.QueryRow(ctx, stmt, args...)

	var im models.InboundMessage

	err := r.scanRow(row, &im)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, models.ErrNotFound
	case err == nil:
		return &im, nil
	default:
		return nil, err
	}
}

func (r *InboundMessageRepo) scanRow(row pgx.Row, m *models.InboundMessage) error {
	err := row.Scan(
		&m.ID,
		&m.Shortcode,
		&m.Status,
		&m.MSISDN,
		&m.Body,
		&m.ProviderID,
		&m.ProviderMessageID,
		&m.NextAttemptAt,
		&m.AttemptCounter,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *InboundMessageRepo) wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameProviderMessageID {
			return models.ErrDuplicateProviderMessageID
		}
	}

	return err
}
