package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type subscriber struct {
	topic       string
	qos         byte
	callback    mqtt.MessageHandler
	timeoutTime time.Duration
}
