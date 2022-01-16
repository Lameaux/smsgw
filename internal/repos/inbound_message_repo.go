package repos

import (
	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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

func (r *InboundMessageRepo) FindByID(shortcode, ID string) (*models.InboundMessage, error) {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := selectInboundMessagesBase + "where shortcode = $1 AND id = $2"

	row := r.db.QueryRow(ctx, stmt, shortcode, ID)

	var im models.InboundMessage

	switch err := scanInboundMessage(row, &im); err {
	case pgx.ErrNoRows:
		return nil, nil
	case nil:
		return &im, nil
	default:
		return nil, err
	}
}

func (r *InboundMessageRepo) FindByShortCode(shortcode string) ([]*models.InboundMessage, error) {
	ctx, cancel := DBQueryContext()
	defer cancel()

	var messages []*models.InboundMessage

	stmt := selectInboundMessagesBase + "where shortcode = $1"

	rows, err := r.db.Query(ctx, stmt, shortcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.InboundMessage
		err = scanInboundMessage(rows, &m)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}

	return messages, nil
}

func scanInboundMessage(row pgx.Row, m *models.InboundMessage) error {
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

func (r *InboundMessageRepo) UpdateStatus(im *models.InboundMessage) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := `update inbound_messages
	set status = $1, updated_at = $2
	where id = $3`

	_, err := r.db.Exec(ctx, stmt,
		im.Status,
		im.UpdatedAt,
		im.ID,
	)

	return err
}

func (r *InboundMessageRepo) wrapError(err error) error {
	if pgerr, ok := err.(*pgconn.PgError); ok {
		if pgerr.ConstraintName == constraintNameProviderMessageID {
			return models.ErrDuplicateProviderMessageID
		}
	}

	return err
}
