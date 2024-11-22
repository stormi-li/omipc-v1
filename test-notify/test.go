package main

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omipc-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omipc := omipc.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	for i := 0; i < 100000; i++ {
		omipc.Notify("channel", strconv.Itoa(i))
		time.Sleep(10*time.Microsecond)
	}
	omipc.Notify("channel", "close")
	c:=redis.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	c.Close()
	omipc.Close()
	select{}
}
