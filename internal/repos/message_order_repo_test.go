package repos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
)

func TestMessageOrderRepo_Save(t *testing.T) {
	r := NewMessageOrderRepo(TestAppConfig.DBPool)
	merchantID := models.NewUUID()

	// New message order
	mo := models.NewSimpleMessageOrder(merchantID, "hello world")
	err := r.Save(mo)
	require.NoError(t, err, "Saving MessageOrder failed")
	assert.NotNil(t, mo.ID, "ID is not set")

	// Duplicate ClientTransactionID
	mo2 := models.NewSimpleMessageOrder(merchantID, "hello world")
	mo2.ClientTransactionID = mo.ClientTransactionID
	err = r.Save(mo2)
	require.ErrorAs(t, err, &models.ErrDuplicateClientTransactionID, "Duplicate MessageOrder created")
}

func TestMessageOrderRepo_FindByID(t *testing.T) {
	r := NewMessageOrderRepo(TestAppConfig.DBPool)

	// Unknown IDs
	_, err := r.FindByID(models.NewUUID())
	assert.Equal(t, models.ErrNotFound, err, "ErrNotFound expected")

	// Populate database
	merchantID := models.NewUUID()
	mo := models.NewSimpleMessageOrder(merchantID, "hello world")
	err = r.Save(mo)
	require.NoError(t, err, "Saving MessageOrder failed")

	// Find order
	found, err := r.FindByID(mo.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "MessageOrder not found")
	assert.Equal(t, mo, found, "MessageOrder is not equal")
}

func TestMessageOrderRepo_FindByMerchantAndID(t *testing.T) {
	r := NewMessageOrderRepo(TestAppConfig.DBPool)

	// Unknown IDs
	_, err := r.FindByMerchantAndID(models.NewUUID(), models.NewUUID())
	assert.Equal(t, models.ErrNotFound, err, "ErrNotFound expected")

	// Populate database
	merchantID := models.NewUUID()
	mo := models.NewSimpleMessageOrder(merchantID, "hello world")
	err = r.Save(mo)
	require.NoError(t, err, "Saving MessageOrder failed")

	// Find order
	found, err := r.FindByMerchantAndID(merchantID, mo.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "MessageOrder not found")
	assert.Equal(t, mo, found, "MessageOrder is not equal")

	// Invalid merchant ID
	_, err = r.FindByMerchantAndID(models.NewUUID(), mo.ID)
	assert.Equal(t, models.ErrNotFound, err, "ErrNotFound expected")
}

func TestMessageOrderRepo_FindByQuery(t *testing.T) {
	r := NewMessageOrderRepo(TestAppConfig.DBPool)

	// Populate database
	mo1 := models.NewSimpleMessageOrder(models.NewUUID(), "hello world")
	err := r.Save(mo1)
	require.NoError(t, err, "Saving MessageOrder failed")

	mo2 := models.NewSimpleMessageOrder(models.NewUUID(), "hello world")
	err = r.Save(mo2)
	require.NoError(t, err, "Saving MessageOrder failed")

	sp := inputs.SearchParams{
		Offset: 0,
		Limit:  10,
	}

	// Find by merchant
	q := inputs.MessageOrderSearchParams{
		SearchParams: &sp,
		MerchantID:   mo1.MerchantID,
	}

	found, err := r.FindByQuery(&q)
	assert.Nil(t, err, "unexpected error")
	assert.Len(t, found, 1, "MessageOrder not found")
	assert.Equal(t, mo1, found[0], "MessageOrder is not equal")

	// Find by merchant and clientTransactionID
	q = inputs.MessageOrderSearchParams{
		SearchParams:        &sp,
		MerchantID:          mo2.MerchantID,
		ClientTransactionID: mo2.ClientTransactionID,
	}

	found, err = r.FindByQuery(&q)
	assert.Nil(t, err, "unexpected error")
	assert.Len(t, found, 1, "MessageOrder not found")
	assert.Equal(t, mo2, found[0], "MessageOrder is not equal")

	// Find by unknown merchant and clientTransactionID
	q = inputs.MessageOrderSearchParams{
		SearchParams:        &sp,
		MerchantID:          mo1.MerchantID,
		ClientTransactionID: mo2.ClientTransactionID,
	}

	found, err = r.FindByQuery(&q)
	assert.Nil(t, err, "unexpected error")
	assert.Len(t, found, 0, "result should be empty")
}
