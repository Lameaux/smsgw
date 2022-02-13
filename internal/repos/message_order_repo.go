package repos

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

type MessageOrderRepo struct {
	db db.Conn
}

const (
	constraintNameClientTransactionID = "message_orders_client_transaction_id"
	tableNameMessageOrders            = "message_orders"
)

func NewMessageOrderRepo(db db.Conn) *MessageOrderRepo {
	return &MessageOrderRepo{db}
}

func (r *MessageOrderRepo) Save(mo *models.MessageOrder) error {
	sb := dbQueryBuilder().Insert(tableNameMessageOrders).
		Columns(
			"merchant_id",
			"sender",
			"body",
			"client_transaction_id",
			"notification_url",
			"created_at",
			"updated_at",
		).
		Values(
			mo.MerchantID,
			mo.Sender,
			mo.Body,
			mo.ClientTransactionID,
			mo.NotificationURL,
			mo.CreatedAt,
			mo.UpdatedAt,
		).
		Suffix(`RETURNING "id"`)

	if err := dbQuerySingle(r.db, &mo.ID, sb); err != nil {
		return r.wrapError(err)
	}

	return nil
}

func (r *MessageOrderRepo) FindByID(id string) (*models.MessageOrder, error) {
	var order models.MessageOrder

	sb := r.selectBase().Where("id = ?", id)
	err := dbQuerySingle(r.db, &order, sb)

	return &order, err
}

func (r *MessageOrderRepo) FindByMerchantAndID(merchantID, id string) (*models.MessageOrder, error) {
	var order models.MessageOrder

	sb := r.selectBase().Where("merchant_id = ? AND id = ?", merchantID, id)
	err := dbQuerySingle(r.db, &order, sb)

	return &order, err
}

func (r *MessageOrderRepo) FindByQuery(q *inputs.MessageOrderSearchParams) ([]*models.MessageOrder, error) {
	sb := r.selectBase().Where("merchant_id = ?", q.MerchantID)

	if q.ClientTransactionID != nil {
		sb = sb.Where("client_transaction_id = ?", q.ClientTransactionID)
	}

	sb = appendSearchParams(q.SearchParams, sb)

	orders := []*models.MessageOrder{}
	err := dbQueryAll(r.db, &orders, sb)

	return orders, err
}

func (r *MessageOrderRepo) selectBase() sq.SelectBuilder {
	return dbQueryBuilder().Select(
		"id",
		"merchant_id",
		"sender",
		"body",
		"client_transaction_id",
		"notification_url",
		"created_at",
		"updated_at",
	).From(tableNameMessageOrders)
}

func (r *MessageOrderRepo) wrapError(err error) error {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerr.ConstraintName == constraintNameClientTransactionID {
			return models.ErrDuplicateClientTransactionID
		}
	}

	return err
}
