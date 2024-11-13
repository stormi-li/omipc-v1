package omipc

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	manager "github.com/stormi-li/omi-v1/omi-manager"
)

type Producer struct {
	configSearcher *manager.Searcher
	channel        string
	messageChan    chan []byte
}

func (producer *Producer) sendMessage() {
	ticker := time.NewTicker(retry_wait_time * 2)
	for {
		address, _ := producer.configSearcher.SearchByLoadBalancing(producer.channel)
		if address == "" {
			time.Sleep(retry_wait_time)
			continue
		}
		conn, err := net.Dial("tcp", address)
		if err != nil {
			time.Sleep(retry_wait_time)
			continue
		}
		for {
			breakLoop := false
			select {
			case <-ticker.C:
				conn.Close()
				breakLoop = true
			case msg := <-producer.messageChan:
				_, err := conn.Write(msg)
				if err != nil {
					producer.messageChan <- msg
					breakLoop = true
				}
			}
			if breakLoop {
				break
			}
		}
	}
}

func (producer *Producer) Publish(message []byte) error {
	//长度前缀协议
	byteMessage := []byte(string(message))
	messageLength := uint32(len(byteMessage))

	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, messageLength)
	count := 0
	for {
		select {
		case producer.messageChan <- append(lengthBuf, byteMessage...):
			return nil
		default:
			if count >= 10 {
				return fmt.Errorf("生产队列已满，请检查消费者状态，或增加队列容量")
			}
			count++
			time.Sleep(retry_wait_time)
		}
	}
}
