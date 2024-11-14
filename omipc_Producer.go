package omipc

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	manager "github.com/stormi-li/omi-v1/omi-manager"
)

type Producer struct {
	configSearcher *manager.Searcher
	channel        string
	conn           net.Conn
	rlock          sync.RWMutex
}

func (producer *Producer) connect() error {
	address, _ := producer.configSearcher.SearchByLoadBalancing(producer.channel)
	if address == "" {
		producer.conn = nil
		return fmt.Errorf("no message queue service was found")
	}
	conn, err := net.Dial("tcp", address)
	if err == nil {
		producer.rlock.Lock()
		if producer.conn != nil {
			producer.conn.Close()
		}
		producer.conn = conn
		producer.rlock.Unlock()
		return nil
	}
	return err
}

func (producer *Producer) Publish(message []byte) error {
	var err error
	retryCount := 0

	//长度前缀协议
	byteMessage := []byte(string(message))
	messageLength := uint32(len(byteMessage))

	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, messageLength)

	for {
		producer.rlock.RLock()
		if producer.conn != nil {
			_, err = producer.conn.Write(append(lengthBuf, byteMessage...))
		} else {
			err = fmt.Errorf("no message queue service was found")
		}
		producer.rlock.RUnlock()
		if err == nil {
			break
		}
		newErr := producer.connect()
		if newErr != nil {
			err = newErr
		}
		time.Sleep(retry_wait_time)
		if retryCount == 10 {
			break
		}
		retryCount++
	}
	return err
}
