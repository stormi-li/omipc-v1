# Omipc 异步通知框架
**作者**: stormi-li  
**Email**: 2785782829@qq.com  

## 简介

**Omipc** 是基于 Redis 的 Pub/Sub 机制的消息通知框架。它具有简单的功能，适合用于广播和监听消息通知任务。开发者可以轻松实现异步通知和消息广播。


## 功能

- **广播消息**：支持将消息广播到指定的 Redis 频道。
- **监听消息**：可以订阅指定的 Redis 频道并监听消息的发布。


## 教程
### 安装
```shell
go get github.com/stormi-li/omipc-v1
```
### 广播消息
```go
package main

import (
	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omipc-v1"
)

func main() {
	omipc := omipc.NewClient(&redis.Options{Addr: "localhost:6379"})
	
	// Notify 用于向指定频道发送消息
	omipc.Notify("channel", "hello world")
}
```
### 监听消息 Listen
```go
package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	omipc "github.com/stormi-li/omipc-v1"
)

func main() {
	c := omipc.NewClient(&redis.Options{Addr: Addr: "localhost:6379"})
	
    // - channel: 订阅的 Redis 频道
    // - timeout: 超时时间，为 0 表示无超时
    // - handFuncs: 可选的处理函数，用于处理收到的消息
    // 返回值：如果收到消息并未超时，返回消息的内容；超时时返回空字符串。
	c.Listen("channel", 0, func(message string) {
		fmt.Println(message)
	})
}
```