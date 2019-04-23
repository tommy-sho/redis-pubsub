package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Client struct {
	client *redis.Pool
}

func NewClient(host string) (*Client, error) {
	p := newPool(host)
	return &Client{p}, nil
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func (c Client) Set(ctx context.Context, key, value string) error {
	conn := c.client.Get()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("redis set error :%v ", err)
	}

	return nil
}

func (c Client) Get(ctx context.Context, key string) (string, error) {
	conn := c.client.Get()
	s, err := redis.String(conn.Do("GET", key))
	fmt.Println(err)
	if err != nil {
		if err == redis.ErrNil {
			return "", err
		}
		return "", fmt.Errorf("redis get error :%v ", err)
	}

	return s, nil
}

func (c Client) MultiGet(ctx context.Context, key []string) ([]string, error) {
	conn := c.client.Get()
	var ts []interface{}
	for _, k := range key {
		ts = append(ts, k)
	}

	values, err := redis.Strings(conn.Do("MGET", ts...))
	fmt.Println(err)
	if err != nil {
		if err == redis.ErrNil {
			return values, err
		}
		return nil, fmt.Errorf("redis get error :%v ", err)
	}

	return values, nil
}
