package message

import (
	corerepos "github.com/Lameaux/core/repos"
	"github.com/Lameaux/smsgw/internal/config"
	"github.com/Lameaux/smsgw/internal/notifications"
	nm "github.com/Lameaux/smsgw/internal/notifications/models"
	"github.com/Lameaux/smsgw/internal/outbound/models"
	"github.com/Lameaux/smsgw/internal/outbound/outputs"
	rg "github.com/Lameaux/smsgw/internal/outbound/repos/group"
	rm "github.com/Lameaux/smsgw/internal/outbound/repos/message"

	im "github.com/Lameaux/smsgw/internal/outbound/inputs/message"
)

type Service struct {
	app *config.App
}

func NewService(app *config.App) *Service {
	return &Service{app}
}

func (s *Service) FindByMerchantAndID(merchantID, id string) (*outputs.MessageView, error) {
	ctx, cancel := corerepos.DBConnContext()
	defer cancel()

	conn, err := s.app.Config.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	message, err := rm.NewRepo(conn).FindByMerchantAndID(merchantID, id)
	if err != nil {
		return nil, err
	}

	messageGroup, err := rg.NewRepo(conn).FindByID(message.MessageGroupID)
	if err != nil {
		return nil, err
	}

	return outputs.NewMessageView(message, messageGroup), nil
}

func (s *Service) AckByProviderAndMessageID(providerID, messageID string) (*models.Message, error) {
	tx, err := corerepos.Begin(s.app.Config.DBPool)
	if err != nil {
		return nil, err
	}

	defer corerepos.Rollback(tx)

	repo := rm.NewRepo(tx)

	message, err := repo.FindByProviderAndMessageID(providerID, messageID)
	if err != nil {
		return nil, err
	}

	messageGroup, err := rg.NewRepo(tx).FindByID(message.MessageGroupID)
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

	if messageGroup.NotificationURL != nil {
		notificationRepo := notifications.NewRepo(tx)
		n := nm.MakeOutboundDeliveryNotification(message)

		err = notificationRepo.Save(n)
		if err != nil {
			return nil, err
		}
	}

	if err := corerepos.Commit(tx); err != nil {
		return nil, err
	}

	return message, err
}

func (s *Service) FindByQuery(p *im.SearchParams) ([]*models.Message, error) {
	ctx, cancel := corerepos.DBConnContext()
	defer cancel()

	conn, err := s.app.Config.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	messages, err := rm.NewRepo(conn).FindByQuery(p)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
