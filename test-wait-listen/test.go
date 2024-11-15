package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omipc-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omipc := omipc.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	key := "channel:channel:channel:channel:channel:channel"
	c := omipc.Listen(key, func(message string) bool {
		fmt.Println(message)
		return message != "close"
	})
	omipc.Notify(key, "hello")
	omipc.Notify(key, "close")
	<-c
}

func wait() {
	omipc := omipc.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	for i := 0; i < 10000; i++ {
		fmt.Println(omipc.Wait("channel", 1*time.Second))
		time.Sleep(1000 * time.Millisecond)
	}
}

func listen() {
	omipc := omipc.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	omipc.Listen("channel", func(message string) bool {
		fmt.Println(message)
		return true
	})
}
