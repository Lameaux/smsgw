package repos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"euromoby.com/smsgw/internal/models"
)

func TestOutboundMessageRepo_Save(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	om := models.NewOutboundMessage(models.NewUUID(), models.NewUUID(), 420123456789)
	err := r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	assert.NotNil(t, om.ID, "ID is not set")
	assert.Equal(t, models.OutboundMessageStatusNew, om.Status, "Invalid status")
}

func TestOutboundMessageRepo_FindByMerchantAndID(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	found, err := r.FindByMerchantAndID(models.NewUUID(), models.NewUUID())
	assert.Nil(t, err, "unexpected error")
	assert.Nil(t, found, "wrong OutboundMessage found")

	merchantID := models.NewUUID()
	om := models.NewOutboundMessage(merchantID, models.NewUUID(), 420123456789)
	err = r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	found, err = r.FindByMerchantAndID(merchantID, om.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "OutboundMessage not found")

	assert.Equal(t, om, found, "OutboundMessage is not equal")
}

func TestOutboundMessageRepo_FindByMerchantAndOrderID(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	res, err := r.FindByMerchantAndOrderID(models.NewUUID(), models.NewUUID())
	require.NoError(t, err, "Finding OutboundMessage failed")
	assert.Equal(t, []*models.OutboundMessage(nil), res, "wrong OutboundMessage found")

	merchantID := models.NewUUID()
	messageOrderID := models.NewUUID()
	om := models.NewOutboundMessage(merchantID, messageOrderID, 420123456789)
	err = r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	found, err := r.FindByMerchantAndOrderID(merchantID, messageOrderID)
	require.NoError(t, err, "Finding OutboundMessage failed")
	assert.NotEmpty(t, found, "OutboundMessage not found")
	assert.Equal(t, om, found[0], "OutboundMessage is not equal")
}
