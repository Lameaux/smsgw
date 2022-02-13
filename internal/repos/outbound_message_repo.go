package repos

import (
	sq "github.com/Masterminds/squirrel"

	"euromoby.com/smsgw/internal/db"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

const (
	tableNameOutboundMessages = "outbound_messages"
)

type OutboundMessageRepo struct {
	db db.Conn
}

func NewOutboundMessageRepo(db db.Conn) *OutboundMessageRepo {
	return &OutboundMessageRepo{db}
}

func (r *OutboundMessageRepo) Save(om *models.OutboundMessage) error {
	sb := dbQueryBuilder().Insert(tableNameOutboundMessages).
		Columns(
			"merchant_id",
			"message_order_id",
			"status",
			"msisdn",
			"next_attempt_at",
			"attempt_counter",
			"created_at",
			"updated_at",
		).
		Values(
			om.MerchantID,
			om.MessageOrderID,
			om.Status,
			om.MSISDN,
			om.NextAttemptAt,
			om.AttemptCounter,
			om.CreatedAt,
			om.UpdatedAt,
		).
		Suffix(`RETURNING "id"`)

	return dbQuerySingle(r.db, &om.ID, sb)
}

func (r *OutboundMessageRepo) Update(om *models.OutboundMessage) error {
	om.UpdatedAt = models.TimeNow()

	sb := dbQueryBuilder().Update(tableNameOutboundMessages).SetMap(
		map[string]interface{}{
			"status":              om.Status,
			"provider_id":         om.ProviderID,
			"provider_response":   om.ProviderResponse,
			"provider_message_id": om.ProviderMessageID,
			"next_attempt_at":     om.NextAttemptAt,
			"attempt_counter":     om.AttemptCounter,
			"updated_at":          om.UpdatedAt,
		},
	).Where("id = ?", om.ID)

	return dbExec(r.db, sb)
}

func (r *OutboundMessageRepo) FindByID(id string) (*models.OutboundMessage, error) {
	var msg models.OutboundMessage

	sb := r.selectBase().Where("id = ?", id)
	err := dbQuerySingle(r.db, &msg, sb)

	return &msg, err
}

func (r *OutboundMessageRepo) FindByMerchantAndID(merchantID, id string) (*models.OutboundMessage, error) {
	var msg models.OutboundMessage

	sb := r.selectBase().
		Where("merchant_id = ?", merchantID).
		Where("id = ?", id)

	err := dbQuerySingle(r.db, &msg, sb)

	return &msg, err
}

func (r *OutboundMessageRepo) FindByMerchantAndOrderID(merchantID, orderID string) ([]*models.OutboundMessage, error) {
	sb := r.selectBase().
		Where("merchant_id = ?", merchantID).
		Where("message_order_id = ?", orderID)

	messages := []*models.OutboundMessage{}
	err := dbQueryAll(r.db, &messages, sb)

	return messages, err
}

func (r *OutboundMessageRepo) FindByProviderAndMessageID(providerID, messageID string) (*models.OutboundMessage, error) {
	var msg models.OutboundMessage

	sb := r.selectBase().
		Where("provider_id = ?", providerID).
		Where("provider_message_id = ?", messageID)

	err := dbQuerySingle(r.db, &msg, sb)

	return &msg, err
}

func (r *OutboundMessageRepo) FindByQuery(q *inputs.OutboundMessageSearchParams) ([]*models.OutboundMessage, error) {
	sb := r.selectBase().Where("merchant_id = ?", q.MerchantID)

	sb = appendMessageParams(q.MessageParams, sb)
	sb = appendSearchParams(q.SearchParams, sb)

	messages := []*models.OutboundMessage{}
	err := dbQueryAll(r.db, &messages, sb)

	return messages, err
}

func (r *OutboundMessageRepo) FindOneForProcessing() (*models.OutboundMessage, error) {
	var msg models.OutboundMessage

	sb := r.selectBase().
		Where("status = ?", models.OutboundMessageStatusNew).
		Where("next_attempt_at < ?", models.TimeNow()).
		Suffix("for update skip locked").
		Limit(1)

	err := dbQuerySingle(r.db, &msg, sb)

	return &msg, err
}

func (r *OutboundMessageRepo) selectBase() sq.SelectBuilder {
	return dbQueryBuilder().Select(
		"id",
		"merchant_id",
		"message_order_id",
		"status",
		"msisdn",
		"provider_id",
		"provider_message_id",
		"provider_response",
		"next_attempt_at",
		"attempt_counter",
		"created_at",
		"updated_at",
	).From(tableNameOutboundMessages)
}
