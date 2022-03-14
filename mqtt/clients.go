package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type client struct {
	*Opts
	conn       mqtt.Client
	isConnect  bool
	subscriber []subscriber
}

func (c *client) Connect() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.Addr, c.Port))
	opts.SetClientID(c.ClientID)
	opts.SetUsername(c.Username)
	opts.SetPassword(c.Password)
	opts.SetOrderMatters(c.Order)
	opts.SetCleanSession(c.CleanSession)
	opts.SetAutoReconnect(c.AutoReconnect)
	opts.SetConnectRetry(c.ConnectRetry)
	opts.SetConnectRetryInterval(c.ConnectRetryInterval)
	opts.OnConnect = func(mc mqtt.Client) {
		c.isConnect = true
		go c.reSubscribe()
		c.OnConnect(mc)
	}
	opts.OnConnectionLost = func(mc mqtt.Client, err error) {
		c.isConnect = false
		c.OnClose(mc, err)
	}

	c.conn = mqtt.NewClient(opts)
}

func (c *client) UnConnect() {
	c.conn.Disconnect(250)
}

func (c *client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler, timeoutTime time.Duration) error {
	token := c.conn.Subscribe(topic, qos, callback)
	token.WaitTimeout(timeoutTime)
	if err := token.Error(); err != nil {
		return err
	}

	c.subscriber = append(c.subscriber, subscriber{
		topic:       topic,
		qos:         qos,
		callback:    callback,
		timeoutTime: timeoutTime,
	})
	return nil
}

func (c *client) reSubscribe() {
	for _, v := range c.subscriber {
		token := c.conn.Subscribe(v.topic, v.qos, v.callback)
		token.WaitTimeout(v.timeoutTime)
		if err := token.Error(); err != nil {

		}
	}
}

func (c *client) Unsubscribe(timeoutTime time.Duration, topics ...string) error {
	token := c.conn.Unsubscribe(topics...)
	token.WaitTimeout(timeoutTime)
	if err := token.Error(); err != nil {
		return err
	}
	return nil
}

func (c *client) Publish(topic string, qos byte, retained bool, payload []byte, timeoutTime time.Duration) error {
	token := c.conn.Publish(topic, qos, retained, payload)
	token.WaitTimeout(timeoutTime)
	if err := token.Error(); err != nil {
		return err
	}
	return nil
}

func (c *client) IsConnected() bool {
	if c.AutoReconnect {
		return c.IsConnected()
	}
	return c.isConnect
}
