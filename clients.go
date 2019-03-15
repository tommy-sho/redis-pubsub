package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Publisher struct{
	redis *redis.Client
	channel string
	pubsub *redis.PubSub
}

func NewPublisher(channel string) *Publisher {
	client := NewRedis()
	return &Publisher{
		redis:client,
		channel:channel,
		pubsub:client.Subscribe(channel),
	}
}

func (p Publisher) SubChannel() <-chan *redis.Message{
	_, err := p.pubsub.Receive()
	if err != nil {
		panic(err)
	}

	return p.pubsub.Channel()
}

func (p Publisher)Close() error {
	err := p.pubsub.Close()
	return err
}


func (p Publisher) Publish(message string) error {
	err := p.redis.Publish(p.channel,message).Err()
	return err
}


type Subscriber struct {
	ch <-chan *redis.Message
}

func (s Subscriber)RecieveMessage() {
	for msg:= range s.ch{
		fmt.Println("recieve: ",msg)
	}
}