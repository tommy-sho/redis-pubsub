package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"time"
)


func NewRedis() *redis.Client{
	client := redis.NewClient(
		&redis.Options{
			Addr : "localhost:6379",
			Password:"",
			DB:0,
		})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}

func main() {

	publisher := NewPublisher("channel1")


	subscriber := Subscriber{
		ch: publisher.SubChannel(),
	}


	go func() {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		err := publisher.Publish(time.Now().String())
		if err != nil {
			log.Fatal(err)
		}
	}
	publisher.Close()
}()

	// Consume messages.
	for msg := range subscriber.ch {
		fmt.Println(msg.Payload)
	}

}