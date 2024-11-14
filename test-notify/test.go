package main

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omipc-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omipc := omipc.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	for i := 0; i < 10; i++ {
		omipc.Notify("channel", strconv.Itoa(i))
	}
	omipc.Notify("channel", "close")
}
