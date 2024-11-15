package omipc

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Client struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (c *Client) Close() {
	c.redisClient.Close()
}

func (c *Client) Listen(channel string, handler func(message string) bool) chan struct{} {
	sub := c.redisClient.Subscribe(c.ctx, channel)
	msgChan := sub.Channel()
	shutdown := make(chan struct{}, 1)
	go func() {
		defer sub.Close()
		for {
			breakLoop := false
			select {
			case msg := <-msgChan:
				if !handler(msg.Payload) {
					breakLoop = true
				}
			case <-shutdown:
				breakLoop = true
			}
			if breakLoop {
				break
			}
		}
		shutdown <- struct{}{}
	}()
	return shutdown
}

func (c *Client) Notify(channel, msg string) {
	c.redisClient.Publish(c.ctx, channel, msg)
}

func (c *Client) Wait(channel string, timeout time.Duration) string {
	sub := c.redisClient.Subscribe(c.ctx, channel)
	defer sub.Close()
	msgChan := sub.Channel()
	return c.wait(msgChan, timeout)
}

func (c *Client) wait(msgChan <-chan *redis.Message, timeout time.Duration) string {
	if timeout == 0 {
		msg := <-msgChan
		return msg.Payload
	}

	timer := time.NewTicker(timeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return ""
	case msg := <-msgChan:
		return msg.Payload
	}
}

func (c *Client) NewLock(lockName string) *Lock {
	return &Lock{
		uuid:        uuid.NewString(),
		lockName:    lockName,
		stop:        make(chan struct{}, 1),
		omipcClient: c,
		redisClient: c.redisClient,
		ctx:         context.Background(),
	}
}
