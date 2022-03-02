package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
)

type Conf struct {
	InsName  string `json:"ins_name"`
	ClientId string `json:"client_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Addr     string `json:"addr"`
	Port     string `json:"port"`
}

type SubscribeOpts struct {
	Topic    string
	Qos      byte
	Callback func(client mqtt.Client, msg mqtt.Message)
}

func InitMqtt(confs []Conf, subOpts map[string][]SubscribeOpts) {
	for k := range confs {
		conf := confs[k]
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s:%s", conf.Addr, conf.Port))
		opts.SetClientID(conf.ClientId)
		opts.SetUsername(conf.Username)
		opts.SetPassword(conf.Password)
		opts.CleanSession = false
		opts.AutoReconnect = true
		opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
			fmt.Println("OnMessage-------:"+conf.InsName, "msg:"+string(msg.Payload()), "topic:"+msg.Topic(), "msgId:"+strconv.Itoa(int(msg.MessageID())))
		})
		opts.OnConnect = func(client mqtt.Client) {
			fmt.Println("OnConnect-------:" + conf.InsName)
		}
		opts.OnConnectionLost = func(client mqtt.Client, err error) {
			fmt.Println("OnClose-------:"+conf.InsName, err)
		}

		registerMqttClient(conf, opts, subOpts[conf.InsName])
	}
}

func registerMqttClient(conf Conf, mqttOpts *mqtt.ClientOptions, opts []SubscribeOpts) {

	c := &Client{
		client: mqtt.NewClient(mqttOpts),
		opts:   opts,
		conf:   conf,
	}

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if err := c.subscribe(); err != nil {
		panic(err)
	}

	mqttClients[conf.InsName] = c
}
