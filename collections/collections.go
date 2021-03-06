package collections

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tokenLen = 44

type ctxKey string

type Collections struct {
	m *mongo.Client

	Sessions *Sessions
	Accounts *Accounts
	Invites  *Invites

	Threads  *Threads
	APIKeys  *APIKeys
	IPNSKeys *IPNSKeys

	Users *Users
}

// NewCollections gets or create store instances for active collections.
func NewCollections(ctx context.Context, uri, dbName string) (*Collections, error) {
	m, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := m.Database(dbName)

	sessions, err := NewSessions(ctx, db)
	if err != nil {
		return nil, err
	}
	accounts, err := NewAccounts(ctx, db)
	if err != nil {
		return nil, err
	}
	invites, err := NewInvites(ctx, db)
	if err != nil {
		return nil, err
	}
	threads, err := NewThreads(ctx, db)
	if err != nil {
		return nil, err
	}
	apikeys, err := NewAPIKeys(ctx, db)
	if err != nil {
		return nil, err
	}
	ipnskeys, err := NewIPNSKeys(ctx, db)
	if err != nil {
		return nil, err
	}
	users, err := NewUsers(ctx, db)
	if err != nil {
		return nil, err
	}

	return &Collections{
		m: m,

		Sessions: sessions,
		Accounts: accounts,
		Invites:  invites,

		Threads:  threads,
		APIKeys:  apikeys,
		IPNSKeys: ipnskeys,

		Users: users,
	}, nil
}

func (c *Collections) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return c.m.Disconnect(ctx)
}
