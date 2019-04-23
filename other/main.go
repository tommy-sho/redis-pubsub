package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

func NewRedis(host string, port string) redis.Conn {
	IP_PORT := fmt.Sprintf("%v:%v", host, port)
	fmt.Println(IP_PORT)
	c, err := redis.Dial("tcp", IP_PORT)
	if err != nil {
		panic(err)
	}

	return c
}

func Set(c redis.Conn, key, value string) error {
	_, err := c.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("redis set error :%v ", err)
	}

	return nil
}

func Exist(c redis.Conn, key string) (bool, error) {
	s, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false, fmt.Errorf("redis exists error :%v ", err)
	}

	return s, nil
}

func Get(c redis.Conn, key string) (string, error) {
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		fmt.Println(err)
	}

	return s, nil
}

func main() {
	c := NewRedis("localhost", "6379")
	defer func() {
		err := c.Close()
		if err != nil {
			log.Fatal("canot close redis connection")
		}
	}()

	s := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

L:
	for s.Scan() {
		n := s.Text()
		switch n {
		case "set":
			var key, value string
			fmt.Print("key: > ")
			if s.Scan() {
				key = s.Text()
			}
			fmt.Print("value: > ")
			if s.Scan() {
				value = s.Text()
			}
			err := Set(c, key, value)
			if err != nil {
				fmt.Println("set error")
			}
			fmt.Println("set ", key, ": ", value)
			fmt.Print("> ")
		case "get":
			var key string
			fmt.Print("key: > ")
			if s.Scan() {
				key = s.Text()
			}

			v, err := Get(c, key)
			if err != nil {
				fmt.Println(err)
			}
			if v == "" {
				fmt.Println("value はありません")
			}
			fmt.Println("value: ", v)
			fmt.Print("> ")
		case "exist":
			var key string
			fmt.Print("key: > ")
			if s.Scan() {
				key = s.Text()
			}

			v, err := Exist(c, key)
			if err != nil {
				fmt.Println(err)
			}
			if v {
				fmt.Println("存在します")
			} else {
				fmt.Println("存在しません")
			}
			fmt.Print("> ")
		case "ex":
			break L
		default:
		}
	}
}
