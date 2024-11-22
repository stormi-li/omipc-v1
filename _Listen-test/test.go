package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omipc-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	c := omipc.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	c.Listen("channel", 0, func(message string) {
		fmt.Println(message)
	})
}
