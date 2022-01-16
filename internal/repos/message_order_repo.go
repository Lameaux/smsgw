package repos

import (
	"fmt"
	"time"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type MessageOrderQuery struct {
	Offset int
	Limit  int

	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time

	ClientTransactionID *string
}

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

func (r *MessageOrderRepo) FindByID(merchantID, ID string) (*models.MessageOrder, error) {
	stmt := selectMessageOrdersBase + "where merchant_id = $1 AND id = $2"
	ctx, cancel := DBQueryContext()
	defer cancel()
	row := r.db.QueryRow(ctx, stmt, merchantID, ID)

	var mo models.MessageOrder
	switch err := r.scanMessageOrder(row, &mo); err {
	case pgx.ErrNoRows:
		return nil, nil
	case nil:
		return &mo, nil
	default:
		return nil, err
	}
}

func (r *MessageOrderRepo) FindByQuery(merchantID string, q *MessageOrderQuery) ([]*models.MessageOrder, error) {
	orders := []*models.MessageOrder{}

	stmt := selectMessageOrdersBase
	args := make([]interface{}, 0)

	args = append(args, merchantID)
	stmt += fmt.Sprintf("where merchant_id = $%d\n", len(args))

	if q.ClientTransactionID != nil {
		args = append(args, q.ClientTransactionID)
		stmt += fmt.Sprintf("and client_transaction_id = $%d\n", len(args))
	}

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

func (r *MessageOrderRepo) wrapError(err error) error {
	if pgerr, ok := err.(*pgconn.PgError); ok {
		if pgerr.ConstraintName == constraintNameClientTransactionID {
			return models.ErrDuplicateProviderMessageID
		}
	}

	return err
}
