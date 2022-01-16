package repos

import (
	"testing"

	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutboundMessageRepo_Save(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	om := models.NewOutboundMessage(utils.NewUUID(), utils.NewUUID(), "420123456789")
	err := r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	assert.NotNil(t, om.ID, "ID is not set")
	assert.Equal(t, models.OutboundMessageStatusNew, om.Status, "Invalid status")
}

func TestOutboundMessageRepo_FindByID(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	found, err := r.FindByID(utils.NewUUID(), utils.NewUUID())
	assert.Nil(t, err, "unexpected error")
	assert.Nil(t, found, "wrong OutboundMessage found")

	merchantID := utils.NewUUID()
	om := models.NewOutboundMessage(merchantID, utils.NewUUID(), "420123456789")
	err = r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	found, err = r.FindByID(merchantID, om.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "OutboundMessage not found")

	assert.Equal(t, om, found, "OutboundMessage is not equal")
}

func TestOutboundMessageRepo_FindByMessageOrderID(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	res, err := r.FindByMessageOrderID(utils.NewUUID(), utils.NewUUID())
	require.NoError(t, err, "Finding OutboundMessage failed")
	assert.Equal(t, []*models.OutboundMessage(nil), res, "wrong OutboundMessage found")

	merchantID := utils.NewUUID()
	messageOrderID := utils.NewUUID()
	om := models.NewOutboundMessage(merchantID, messageOrderID, "420123456789")
	err = r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	found, err := r.FindByMessageOrderID(merchantID, messageOrderID)
	require.NoError(t, err, "Finding OutboundMessage failed")
	assert.NotEmpty(t, found, "OutboundMessage not found")
	assert.Equal(t, om, found[0], "OutboundMessage is not equal")
}
