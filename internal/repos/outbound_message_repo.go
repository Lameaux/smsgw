package repos

import (
	"fmt"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
	"github.com/jackc/pgx/v4"
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

func (r *OutboundMessageRepo) FindByID(merchantID, id string) (*models.OutboundMessage, error) {
	stmt := selectOutboundMessagesBase + "where merchant_id = $1 AND id = $2"

	ctx, cancel := DBQueryContext()
	defer cancel()
	row := r.db.QueryRow(ctx, stmt, merchantID, id)

	var om models.OutboundMessage
	switch err := r.scanRow(row, &om); err {
	case pgx.ErrNoRows:
		return nil, nil
	case nil:
		return &om, nil
	default:
		return nil, err
	}
}

func (r *OutboundMessageRepo) FindByMessageOrderID(merchantID, id string) ([]*models.OutboundMessage, error) {
	var messages []*models.OutboundMessage

	stmt := selectOutboundMessagesBase + "where merchant_id = $1 AND message_order_id = $2"

	ctx, cancel := DBQueryContext()
	defer cancel()
	rows, err := r.db.Query(ctx, stmt, merchantID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.OutboundMessage
		err = r.scanRow(rows, &m)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}

	return messages, nil
}

func (r *OutboundMessageRepo) FindByQuery(q *inputs.OutboundMessageSearchParams) ([]*models.OutboundMessage, error) {
	messages := []*models.OutboundMessage{}

	stmt := selectOutboundMessagesBase
	args := make([]interface{}, 0)

	args = append(args, q.MerchantID)
	stmt += fmt.Sprintf("where merchant_id = $%d\n", len(args))

	stmt, args = appendMessageParams(q.MessageParams, stmt, args)
	stmt, args = appendSearchParams(q.SearchParams, stmt, args)

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

func (r *OutboundMessageRepo) FindOneForProcessing() (*models.OutboundMessage, error) {
	stmt := selectOutboundMessagesBase + `
	where status = $1
	and next_attempt_at < now()
	for update skip locked
 	limit 1
	`
	ctx, cancel := DBQueryContext()
	defer cancel()
	row := r.db.QueryRow(ctx, stmt, models.OutboundMessageStatusNew)

	var om models.OutboundMessage
	switch err := r.scanRow(row, &om); err {
	case pgx.ErrNoRows:
		return nil, nil
	case nil:
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

func (r *OutboundMessageRepo) UpdateStatus(om *models.OutboundMessage) error {
	ctx, cancel := DBQueryContext()
	defer cancel()

	stmt := `update outbound_messages
	set status = $1,
	provider_id = $2, provider_response = $3, provider_message_id = $4,
	updated_at = $5
	where id = $6`

	_, err := r.db.Exec(ctx, stmt,
		om.Status,
		om.ProviderID,
		om.ProviderResponse,
		om.ProviderMessageID,
		om.UpdatedAt,
		om.ID,
	)

	return err
}
