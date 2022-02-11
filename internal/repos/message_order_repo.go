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

type MessageOrderRepo struct {
	db db.Conn
}

const (
	selectMessageOrdersBase = `select
	id, merchant_id, sender, body, client_transaction_id,
	notification_url,
	created_at, updated_at
	from message_orders
	`
	constraintNameClientTransactionID = "message_orders_client_transaction_id"
)

func NewMessageOrderRepo(db db.Conn) *MessageOrderRepo {
	return &MessageOrderRepo{db}
}

func (r *MessageOrderRepo) Save(mo *models.MessageOrder) error {
	var insertedID string

	stmt := `insert into message_orders (
		merchant_id, sender, body, client_transaction_id,
		notification_url, created_at, updated_at
	) values ($1, $2, $3, $4, $5, $6, $7)
	returning id
	`

	ctx, cancel := DBQueryContext()
	defer cancel()

	err := r.db.QueryRow(ctx, stmt,
		mo.MerchantID,
		mo.Sender,
		mo.Body,
		mo.ClientTransactionID,
		mo.NotificationURL,
		mo.CreatedAt,
		mo.UpdatedAt,
	).Scan(&insertedID)
	if err != nil {
		return r.wrapError(err)
	}

	mo.ID = insertedID

	return nil
}

func (r *MessageOrderRepo) FindByID(id string) (*models.MessageOrder, error) {
	stmt := selectMessageOrdersBase + "where id = $1"

	return r.querySingle(stmt, id)
}

func (r *MessageOrderRepo) FindByMerchantAndID(merchantID, id string) (*models.MessageOrder, error) {
	stmt := selectMessageOrdersBase + "where merchant_id = $1 AND id = $2"

	return r.querySingle(stmt, merchantID, id)
}

func (r *MessageOrderRepo) FindByQuery(q *inputs.MessageOrderSearchParams) ([]*models.MessageOrder, error) {
	stmt := selectMessageOrdersBase
	args := make([]interface{}, 0)

	args = append(args, q.MerchantID)
	stmt += fmt.Sprintf("where merchant_id = $%d\n", len(args))

	if q.ClientTransactionID != nil {
		args = append(args, q.ClientTransactionID)
		stmt += fmt.Sprintf("and client_transaction_id = $%d\n", len(args))
	}

	stmt, args = appendSearchParams(q.SearchParams, stmt, args)

	return r.query(stmt, args...)
}

func (r *MessageOrderRepo) query(stmt string, args ...interface{}) ([]*models.MessageOrder, error) {
	orders := []*models.MessageOrder{}

	ctx, cancel := DBQueryContext()
	defer cancel()

	rows, err := r.db.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var mo models.MessageOrder

		err = r.scanMessageOrder(rows, &mo)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &mo)
	}

	return orders, nil
}

func (r *MessageOrderRepo) querySingle(stmt string, args ...interface{}) (*models.MessageOrder, error) {
	ctx, cancel := DBQueryContext()
	defer cancel()

	row := r.db.QueryRow(ctx, stmt, args...)

	var mo models.MessageOrder

	err := r.scanMessageOrder(row, &mo)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, models.ErrNotFound
	case err == nil:
		return &mo, nil
	default:
		return nil, err
	}
}

func (r *MessageOrderRepo) scanMessageOrder(row pgx.Row, mo *models.MessageOrder) error {
	return row.Scan(
		&mo.ID,
		&mo.MerchantID,
		&mo.Sender,
		&mo.Body,
		&mo.ClientTransactionID,
		&mo.NotificationURL,
		&mo.CreatedAt,
		&mo.UpdatedAt,
	)
}

func (r *MessageOrderRepo) wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameClientTransactionID {
			return models.ErrDuplicateProviderMessageID
		}
	}

	return err
}
