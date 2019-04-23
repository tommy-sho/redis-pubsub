package tweetreader

import (
	"context"
)

type Action struct {
	IsDidFired     bool
	IsDidRetweeted bool
	IsDidReply     bool
}

type ActionRepository interface {
	Set(ctx context.Context, accountID string, tweetID string, action *Action) error
	Get(ctx context.Context, accountID string, tweetID string) (*Action, error)
	GetMulti(ctx context.Context, accountID string, tweetIDs []string) ([]*Action, error)
}
