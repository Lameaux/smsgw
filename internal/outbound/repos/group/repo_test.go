package group

import (
	"github.com/Lameaux/smsgw/internal/outbound/inputs/group"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Lameaux/smsgw/internal/outbound/models"

	coremodels "github.com/Lameaux/core/models"
	commoninputs "github.com/Lameaux/smsgw/internal/common/inputs"
)

func TestRepo_Save(t *testing.T) {
	r := NewRepo(TestApp.Config.DBPool)
	merchantID := coremodels.NewUUID()

	// New message group
	mo := models.NewMessageGroup(merchantID, "hello world")
	err := r.Save(mo)
	require.NoError(t, err, "Saving MessageGroup failed")
	assert.NotNil(t, mo.ID, "ID is not set")

	// Duplicate ClientTransactionID
	mo2 := models.NewMessageGroup(merchantID, "hello world")
	mo2.ClientTransactionID = mo.ClientTransactionID
	err = r.Save(mo2)
	require.ErrorAs(t, err, &models.ErrDuplicateClientTransactionID, "Duplicate MessagGroup created")
}

func TestRepo_FindByID(t *testing.T) {
	r := NewRepo(TestApp.Config.DBPool)

	// Unknown IDs
	_, err := r.FindByID(coremodels.NewUUID())
	assert.Equal(t, coremodels.ErrNotFound, err, "ErrNotFound expected")

	// Populate database
	merchantID := coremodels.NewUUID()
	mo := models.NewMessageGroup(merchantID, "hello world")
	err = r.Save(mo)
	require.NoError(t, err, "Saving MessageGroup failed")

	// Find group
	found, err := r.FindByID(mo.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "MessageGroup not found")
	assert.Equal(t, mo, found, "MessageGroup is not equal")
}

func TestRepo_FindByMerchantAndID(t *testing.T) {
	r := NewRepo(TestApp.Config.DBPool)

	// Unknown IDs
	_, err := r.FindByMerchantAndID(coremodels.NewUUID(), coremodels.NewUUID())
	assert.Equal(t, coremodels.ErrNotFound, err, "ErrNotFound expected")

	// Populate database
	merchantID := coremodels.NewUUID()
	mo := models.NewMessageGroup(merchantID, "hello world")
	err = r.Save(mo)
	require.NoError(t, err, "Saving MessageGroup failed")

	// Find group
	found, err := r.FindByMerchantAndID(merchantID, mo.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "MessageGroup not found")
	assert.Equal(t, mo, found, "MessageGroup is not equal")

	// Invalid merchant ID
	_, err = r.FindByMerchantAndID(coremodels.NewUUID(), mo.ID)
	assert.Equal(t, coremodels.ErrNotFound, err, "ErrNotFound expected")
}

func TestRepo_FindByQuery(t *testing.T) {
	r := NewRepo(TestApp.Config.DBPool)

	// Populate database
	mo1 := models.NewMessageGroup(coremodels.NewUUID(), "hello world")
	err := r.Save(mo1)
	require.NoError(t, err, "Saving MessageGroup failed")

	mo2 := models.NewMessageGroup(coremodels.NewUUID(), "hello world")
	err = r.Save(mo2)
	require.NoError(t, err, "Saving MessageGroup failed")

	sp := commoninputs.SearchParams{
		Offset: 0,
		Limit:  10,
	}

	// Find by merchant
	q := group.SearchParams{
		SearchParams: &sp,
		MerchantID:   mo1.MerchantID,
	}

	found, err := r.FindByQuery(&q)
	assert.Nil(t, err, "unexpected error")
	assert.Len(t, found, 1, "MessageGroup not found")
	assert.Equal(t, mo1, found[0], "MessageGroup is not equal")

	// Find by merchant and clientTransactionID
	q = group.SearchParams{
		SearchParams:        &sp,
		MerchantID:          mo2.MerchantID,
		ClientTransactionID: mo2.ClientTransactionID,
	}

	found, err = r.FindByQuery(&q)
	assert.Nil(t, err, "unexpected error")
	assert.Len(t, found, 1, "MessageGroup not found")
	assert.Equal(t, mo2, found[0], "MessageGroup is not equal")

	// Find by unknown merchant and clientTransactionID
	q = group.SearchParams{
		SearchParams:        &sp,
		MerchantID:          mo1.MerchantID,
		ClientTransactionID: mo2.ClientTransactionID,
	}

	found, err = r.FindByQuery(&q)
	assert.Nil(t, err, "unexpected error")
	assert.Len(t, found, 0, "result should be empty")
}
