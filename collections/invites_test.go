package collections_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "github.com/textileio/textile/collections"
)

func TestInvites_Create(t *testing.T) {
	db := newDB(t)
	col, err := NewInvites(context.Background(), db)
	require.Nil(t, err)

	_, from, err := crypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)
	created, err := col.Create(context.Background(), from, "myorg", "jane@doe.com")
	require.Nil(t, err)
	assert.True(t, created.ExpiresAt.After(time.Now()))
}

func TestInvites_Get(t *testing.T) {
	db := newDB(t)
	col, err := NewInvites(context.Background(), db)
	require.Nil(t, err)

	_, from, err := crypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)
	created, err := col.Create(context.Background(), from, "myorg", "jane@doe.com")
	require.Nil(t, err)

	got, err := col.Get(context.Background(), created.Token)
	require.Nil(t, err)
	assert.Equal(t, created.Token, got.Token)
}

func TestInvites_ListByEmail(t *testing.T) {
	db := newDB(t)
	col, err := NewInvites(context.Background(), db)
	require.Nil(t, err)

	list, err := col.ListByEmail(context.Background(), "jane@doe.com")
	require.Nil(t, err)
	require.Empty(t, list)

	_, from, err := crypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)
	created, err := col.Create(context.Background(), from, "myorg", "jane@doe.com")
	require.Nil(t, err)

	list, err = col.ListByEmail(context.Background(), "jane@doe.com")
	require.Nil(t, err)
	require.Equal(t, 1, len(list))
	require.Equal(t, created.Token, list[0].Token)
}

func TestInvites_Accept(t *testing.T) {
	db := newDB(t)
	col, err := NewInvites(context.Background(), db)
	require.Nil(t, err)

	_, from, err := crypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)
	created, err := col.Create(context.Background(), from, "myorg", "jane@doe.com")
	require.Nil(t, err)
	assert.False(t, created.Accepted)

	err = col.Accept(context.Background(), created.Token)
	require.Nil(t, err)
	got, err := col.Get(context.Background(), created.Token)
	require.Nil(t, err)
	assert.True(t, got.Accepted)
}

func TestInvites_Delete(t *testing.T) {
	db := newDB(t)
	col, err := NewInvites(context.Background(), db)
	require.Nil(t, err)

	_, from, err := crypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)
	created, err := col.Create(context.Background(), from, "myorg", "jane@doe.com")
	require.Nil(t, err)

	err = col.Delete(context.Background(), created.Token)
	require.Nil(t, err)
	_, err = col.Get(context.Background(), created.Token)
	require.NotNil(t, err)
}

func TestInvites_DeleteByFrom(t *testing.T) {
	db := newDB(t)
	col, err := NewInvites(context.Background(), db)
	require.Nil(t, err)

	_, from, err := crypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)
	created, err := col.Create(context.Background(), from, "myorg", "jane@doe.com")
	require.Nil(t, err)

	err = col.DeleteByFrom(context.Background(), created.From)
	require.Nil(t, err)
	_, err = col.Get(context.Background(), created.Token)
	require.NotNil(t, err)
}

func TestInvites_DeleteByOrg(t *testing.T) {
	db := newDB(t)
	col, err := NewInvites(context.Background(), db)
	require.Nil(t, err)

	_, from, err := crypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)
	created, err := col.Create(context.Background(), from, "myorg", "jane@doe.com")
	require.Nil(t, err)

	err = col.DeleteByOrg(context.Background(), created.Org)
	require.Nil(t, err)
	_, err = col.Get(context.Background(), created.Token)
	require.NotNil(t, err)
}

func TestInvites_DeleteByFromAndOrg(t *testing.T) {
	db := newDB(t)
	col, err := NewInvites(context.Background(), db)
	require.Nil(t, err)

	_, from, err := crypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)
	created, err := col.Create(context.Background(), from, "myorg", "jane@doe.com")
	require.Nil(t, err)

	err = col.DeleteByFromAndOrg(context.Background(), from, created.Org)
	require.Nil(t, err)
	_, err = col.Get(context.Background(), created.Token)
	require.NotNil(t, err)
}
