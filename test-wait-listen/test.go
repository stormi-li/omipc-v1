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
	listen()
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
	listener := omipc.Listen("channel", func(message string) bool {
		fmt.Println(message)
		if message == "close" {
			return false
		}
		return true
	})
	go func(){
		time.Sleep(1*time.Second)
		listener.Close()
	}()
	listener.Wait()
	fmt.Println("closed")
}
