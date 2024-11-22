package omipc

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (c *Client) Close() {
	c.redisClient.Close()
}

type Listener struct {
	close chan struct{}
	wait  chan struct{}
}

func (l *Listener) Close() {
	l.close <- struct{}{}
}
func (l *Listener) Wait() {
	<-l.wait
}

func (c *Client) Listen(channel string, handler func(message string) bool) *Listener {
	sub := c.redisClient.Subscribe(c.ctx, channel)
	msgChan := sub.Channel()
	close := make(chan struct{}, 1)
	wait := make(chan struct{}, 1)
	go func() {
		defer sub.Close()
		for {
			breakLoop := false
			select {
			case msg := <-msgChan:
				if !handler(msg.Payload) {
					breakLoop = true
				}
			case <-close:
				breakLoop = true
			}
			if breakLoop {
				break
			}
		}
		wait <- struct{}{}
	}()
	return &Listener{
		close: close,
		wait:  wait,
	}
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
