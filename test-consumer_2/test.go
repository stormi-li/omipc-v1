package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
)

func main() {
	consumer()
}

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func consumer() {
	c := omipc.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	consumer := c.NewConsumer("consumer_1", 2, 100000)
	consumer.ListenAndConsume("118.25.196.166:5556", func(message []byte) {
		fmt.Println(string(message))
	})
}
