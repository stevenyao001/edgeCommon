package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var mqttClients map[string]*Client

func init() {
	mqttClients = make(map[string]*Client)
}

type Client struct {
	client mqtt.Client
	opts   []SubscribeOpts
	conf   Conf
}

func (mc *Client) Publish(topic string, msg Msg, qos byte, retained bool) (int64, error) {
	if msg.MsgId == 0 {
		msg.MsgId = time.Now().UnixNano()
	}

	a, _ := json.Marshal(msg)
	token := mc.client.Publish(topic, qos, retained, a)
	_ = token.Wait()
	if err := token.Error(); err != nil {
		return 0, err
	}
	return msg.MsgId, nil
}

func GetClient(insName string) (mc *Client) {
	return mqttClients[insName]
}

func (mc *Client) subscribe() error {

	var token mqtt.Token
	var err error

	for _, opt := range mc.opts {
		token = mc.client.Subscribe(opt.Topic, opt.Qos, opt.Callback)
		_ = token.Wait()
		if err = token.Error(); err != nil {
			err = fmt.Errorf("订阅topic报错:[%s]%s", opt.Topic, err.Error())
			return err
		}
	}

	return nil
}
