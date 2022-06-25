package inbound

import (
	corerepos "github.com/Lameaux/core/repos"
	"github.com/Lameaux/smsgw/internal/config"
	"github.com/Lameaux/smsgw/internal/users"

	"github.com/Lameaux/smsgw/internal/inbound/models"

	"github.com/Lameaux/smsgw/internal/notifications"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
)

type Service struct {
	app *config.App
	us  users.Service
}

func NewService(app *config.App, us users.Service) *Service {
	return &Service{app, us}
}

func (s *Service) FindMerchantByShortcode(shortcode string) (string, error) {
	return s.us.FindMerchantByShortcode(shortcode)
}

func (s *Service) FindByMerchantAndID(merchantID, id string) (*models.Message, error) {
	ctx, cancel := corerepos.DBConnContext()
	defer cancel()

	conn, err := s.app.Config.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	return NewRepo(conn).FindByMerchantAndID(merchantID, id)
}

func (s *Service) SaveMessage(m *models.Message) error {
	ctx, done := corerepos.DBConnContext()
	defer done()

	conn, err := s.app.Config.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	return NewRepo(conn).Save(m)
}

func (s *Service) AckByMerchantAndID(merchantID, id string) (*models.Message, error) {
	tx, err := corerepos.Begin(s.app.Config.DBPool)
	if err != nil {
		return nil, err
	}

	defer corerepos.Rollback(tx)

	repo := NewRepo(tx)

	message, err := repo.FindByMerchantAndID(merchantID, id)
	if err != nil {
		return nil, err
	}

	if message.Status == models.MessageStatusDelivered {
		return nil, models.ErrAlreadyAcked
	}

	message.Status = models.MessageStatusDelivered

	if err := repo.Update(message); err != nil {
		return nil, err
	}

	n := nm.MakeInboundDeliveryNotification(message)

	if err := notifications.NewRepo(tx).Save(n); err != nil {
		return nil, err
	}

	if err := corerepos.Commit(tx); err != nil {
		return nil, err
	}

	return message, err
}

func (s *Service) FindByQuery(p *SearchParams) ([]*models.Message, error) {
	ctx, cancel := corerepos.DBConnContext()
	defer cancel()

	conn, err := s.app.Config.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	messages, err := NewRepo(conn).FindByQuery(p)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
