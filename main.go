package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/tommy-sho/redis-pubsub/tweetreader"

	"github.com/tommy-sho/redis-pubsub/infrastructure"
	"github.com/tommy-sho/redis-pubsub/redis"
)

var (
	accountID = "10000"
	tweetID   = "20000"
)

type Req struct {
	accountID string
	tweetID   string
}

func main() {
	client, err := redis.NewClient("localhost:6379")
	if err != nil {
		panic(err)
	}

	acRep := infrastructure.NewActionRepository(client)
	ctx := context.Background()
	a := &tweetreader.Action{
		IsDidFired:     true,
		IsDidReply:     true,
		IsDidRetweeted: true,
	}

	s := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

L:
	for s.Scan() {
		n := s.Text()
		switch n {
		case "set":
			fmt.Print("accountID: > ")
			if s.Scan() {
				accountID = s.Text()
			}
			fmt.Print("tweetID: > ")
			if s.Scan() {
				tweetID = s.Text()
			}
			err = acRep.Set(ctx, accountID, tweetID, a)
			if err != nil {
				panic(err)
			}
			fmt.Print("> ")
		case "get":
			fmt.Print("accountID: > ")
			if s.Scan() {
				accountID = s.Text()
			}
			fmt.Print("tweetID: > ")
			if s.Scan() {
				tweetID = s.Text()
			}
			a, err = acRep.Get(ctx, accountID, tweetID)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%+v\n", a)
			fmt.Print("> ")
		case "mget":
			var tweetIDs []string
			fmt.Print("accountID: > ")
			if s.Scan() {
				accountID = s.Text()
			}
			for {
				fmt.Print("tweetID: > ")
				if s.Scan() {
					tweetID = s.Text()
				}
				if tweetID == "ex" {
					break
				}
				tweetIDs = append(tweetIDs, tweetID)
			}
			as, err := acRep.GetMulti(ctx, accountID, tweetIDs)
			if err != nil {
				panic(err)
			}
			for _, t := range as {
				fmt.Printf("%+v\n", t)
			}

			fmt.Print("> ")

		case "ex":
			break L
		default:
		}
	}
}
