package main

import (
	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omipc-v1"
)

func main() {
	c := omipc.NewClient(&redis.Options{Addr: "118.25.196.166:3934", Password: "12982397StrongPassw0rd"})
	c.Notify("channel", "hello world")
}
