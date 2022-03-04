package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Opts struct {
	ClientID             string
	Addr                 string
	Port                 int
	Username             string
	Password             string
	Order                bool
	CleanSession         bool
	AutoReconnect        bool
	ConnectRetryInterval time.Duration
	ConnectRetry         bool
	OnConnect            mqtt.OnConnectHandler
	OnClose              mqtt.ConnectionLostHandler
}
