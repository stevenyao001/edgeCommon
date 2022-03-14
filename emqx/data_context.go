// Package emqx group
package emqx

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type DataContest interface {
	SubData()
	PubData()
}

type EMQXGroup struct {
	Clients []mqtt.Client
}

// PubData 发布
func (group *EMQXGroup) PubData(topic string, qos byte, data interface{}) error {
	for _, v := range group.Clients {
		token := v.Publish(topic, qos, false, data)
		if !token.WaitTimeout(time.Second) && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}

// SubData 订阅
func (group *EMQXGroup) SubData(topic string, qos byte, callback func(client mqtt.Client, message mqtt.Message)) error {
	for _, v := range group.Clients {
		token := v.Subscribe(topic, qos, callback)
		if !token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}
