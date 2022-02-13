package repos

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

const (
	selectOutboundMessagesBase = `select
	id, merchant_id, message_order_id, status, msisdn,
	provider_id, provider_message_id, provider_response,
	next_attempt_at, attempt_counter,
	created_at, updated_at
	from outbound_messages
	`
)

type OutboundMessageRepo struct {
	db db.Conn
}

func NewOutboundMessageRepo(db db.Conn) *OutboundMessageRepo {
	return &OutboundMessageRepo{db}
}

func (r *OutboundMessageRepo) Save(om *models.OutboundMessage) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := `insert into outbound_messages (
		merchant_id, message_order_id, status, msisdn,
		next_attempt_at, attempt_counter,
		created_at, updated_at
	)
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	returning id
	`

	var insertedID string

	err := r.db.QueryRow(ctx, stmt,
		om.MerchantID,
		om.MessageOrderID,
		om.Status,
		om.MSISDN,
		om.NextAttemptAt,
		om.AttemptCounter,
		om.CreatedAt,
		om.UpdatedAt,
	).Scan(&insertedID)
	if err != nil {
		return err
	}

	om.ID = insertedID

	return nil
}

func (r *OutboundMessageRepo) Update(om *models.OutboundMessage) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := `update outbound_messages
	set status = $1,
	provider_id = $2, provider_response = $3, provider_message_id = $4,
	next_attempt_at = $5, attempt_counter = $6,
	updated_at = $7
	where id = $8`

	om.UpdatedAt = models.TimeNow()

	_, err := r.db.Exec(ctx, stmt,
		om.Status,
		om.ProviderID,
		om.ProviderResponse,
		om.ProviderMessageID,
		om.NextAttemptAt,
		om.AttemptCounter,
		om.UpdatedAt,
		om.ID,
	)

	return err
}

func (r *OutboundMessageRepo) FindByID(id string) (*models.OutboundMessage, error) {
	stmt := selectOutboundMessagesBase + "where id = $1"

	return r.querySingle(stmt, id)
}

func (r *OutboundMessageRepo) FindByMerchantAndID(merchantID, id string) (*models.OutboundMessage, error) {
	stmt := selectOutboundMessagesBase + "where merchant_id = $1 AND id = $2"

	return r.querySingle(stmt, merchantID, id)
}

func (r *OutboundMessageRepo) FindByMerchantAndOrderID(merchantID, orderID string) ([]*models.OutboundMessage, error) {
	stmt := selectOutboundMessagesBase + "where merchant_id = $1 AND message_order_id = $2"

	return r.query(stmt, merchantID, orderID)
}

func (r *OutboundMessageRepo) FindByProviderAndMessageID(providerID, messageID string) (*models.OutboundMessage, error) {
	stmt := selectOutboundMessagesBase + "where provider_id = $1 AND provider_message_id = $2"

	return r.querySingle(stmt, providerID, messageID)
}

func (r *OutboundMessageRepo) FindByQuery(q *inputs.OutboundMessageSearchParams) ([]*models.OutboundMessage, error) {
	stmt := selectOutboundMessagesBase
	args := make([]interface{}, 0)

	args = append(args, q.MerchantID)
	stmt += fmt.Sprintf("where merchant_id = $%d\n", len(args))

	stmt, args = appendMessageParams(q.MessageParams, stmt, args)
	// stmt, args = appendSearchParams(q.SearchParams, stmt, args)

	return r.query(stmt, args...)
}

func (r *OutboundMessageRepo) FindOneForProcessing() (*models.OutboundMessage, error) {
	stmt := selectOutboundMessagesBase + `
	where status = $1
	and next_attempt_at < $2
	for update skip locked
 	limit 1
	`

	return r.querySingle(stmt, models.OutboundMessageStatusNew, models.TimeNow())
}

func (r *OutboundMessageRepo) query(stmt string, args ...interface{}) ([]*models.OutboundMessage, error) {
	messages := []*models.OutboundMessage{}

	ctx, cancel := DBQueryContext()
	defer cancel()

	rows, err := r.db.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var om models.OutboundMessage

		err = r.scanRow(rows, &om)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &om)
	}

	return messages, nil
}

func (r *OutboundMessageRepo) querySingle(stmt string, args ...interface{}) (*models.OutboundMessage, error) {
	ctx, cancel := DBQueryContext()
	defer cancel()

	row := r.db.QueryRow(ctx, stmt, args...)

	var om models.OutboundMessage

	err := r.scanRow(row, &om)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, models.ErrNotFound
	case err == nil:
		return &om, nil
	default:
		return nil, err
	}
}

func (r *OutboundMessageRepo) scanRow(row pgx.Row, m *models.OutboundMessage) error {
	return row.Scan(
		&m.ID,
		&m.MerchantID,
		&m.MessageOrderID,
		&m.Status,
		&m.MSISDN,
		&m.ProviderID,
		&m.ProviderMessageID,
		&m.ProviderResponse,
		&m.NextAttemptAt,
		&m.AttemptCounter,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
}
