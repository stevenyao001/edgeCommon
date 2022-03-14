package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Clients interface {
	Connect()
	UnConnect()
	Subscribe(topic string, qos byte, callback mqtt.MessageHandler, timeoutTime time.Duration) error
	Unsubscribe(timeoutTime time.Duration, topics ...string) error
	Publish(topic string, qos byte, retained bool, payload []byte, timeoutTime time.Duration) error
	IsConnected() bool
}

func NewClient(opts *Opts) Clients {
	return &client{
		Opts:       opts,
		conn:       nil,
		isConnect:  false,
		subscriber: make([]subscriber, 0),
	}
}
