package ds

import (
	"github.com/asiainfoLDP/datahub_commons/mq"
	"sync"
	"time"
)

type Result struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type MQ struct {
	Mutex        sync.Mutex
	MessageQueue mq.MessageQueue
}

type MesssageListener struct {
	Topic string
	MQ    *MQ

	MessageHandler func(topic string, key, value []byte) error

	closed chan struct{}
	acked  chan struct{} // acknowledged message received
}

type DnsEntry struct {
	Ip   string
	Port string
}

type AlarmEvent struct {
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Send_time time.Time `json:"sendTime"`
}
