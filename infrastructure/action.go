package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/tommy-sho/redis-pubsub/redis"
	"github.com/tommy-sho/redis-pubsub/tweetreader"
)

func NewActionRepository(c *redis.Client) tweetreader.ActionRepository {
	return &ActionRepository{c}
}

type ActionRepository struct {
	client *redis.Client
}

func (a *ActionRepository) Set(ctx context.Context, accountID string, tweetID string, action *tweetreader.Action) error {
	key := NewKey(accountID, tweetID)
	value, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("action repository error :%v ", err)
	}

	err = a.client.Set(key, string(value))
	if err != nil {
		return fmt.Errorf("action repository error : %v ", err)
	}

	return nil
}

func (a *ActionRepository) Get(ctx context.Context, accountID string, tweetID string) (*tweetreader.Action, error) {
	s, err := a.client.Get(NewKey(accountID, tweetID))
	if err != nil {
		return &tweetreader.Action{}, fmt.Errorf("action repository error :%v ", err)
	}

	var res *tweetreader.Action
	err = json.Unmarshal([]byte(s), res)
	if err != nil {
		return &tweetreader.Action{}, fmt.Errorf("action repository: json unmershal error :%v ", err)
	}

	return res, nil
}

func (a *ActionRepository) GetMulti(ctx context.Context, accountID string, tweetIDs []string) ([]*tweetreader.Action, error) {
	res := make([]*tweetreader.Action, len(tweetIDs))
	as := make([]string, len(tweetIDs))
	for i, t := range tweetIDs {
		key := NewKey(accountID, t)
		action, err := a.client.Get(key)
		if err != redigo.ErrNil {
			return nil, fmt.Errorf("action repository error : %v ", err)
		}

		as[i] = action
	}

	for i, t := range as {
		var p *tweetreader.Action
		err := json.Unmarshal([]byte(t), p)
		if err != nil {
			return res, fmt.Errorf("action repository: json unmershal error :%v ", err)
		}

		res[i] = p
	}

	return res, nil
}

func NewKey(accountID, tweetID string) string {
	return fmt.Sprintf("%v/%v", accountID, tweetID)
}
