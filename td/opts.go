package td

type CollectorConf struct {
	InsName      string `json:"ins_name"`
	Driver       string `json:"driver"`
	Network      string `json:"network"`
	Addr         string `json:"addr"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Db           string `json:"db"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxIdleTime  int    `json:"max_idle_time"`
	MaxLifeTime  int    `json:"max_life_time"`
	MaxOpenConns int    `json:"max_open_conns"`
}

func Default(conf CollectorConf) CollectorConf {
	opts := new()
	opts.InsName = conf.InsName
	opts.Driver = conf.Driver
	opts.Network = conf.Network
	opts.Addr = conf.Addr
	opts.Port = conf.Port
	opts.Username = conf.Username
	opts.Password = conf.Password
	opts.Db = conf.Db
	opts.MaxIdleConns = conf.MaxIdleConns
	opts.MaxIdleTime = conf.MaxIdleTime
	opts.MaxLifeTime = conf.MaxLifeTime
	opts.MaxOpenConns = conf.MaxOpenConns
	return opts
}

func new() CollectorConf {
	return CollectorConf{
		InsName:      "",
		Driver:       "taosRestful",
		Network:      "http",
		Addr:         "127.0.0.1",
		Port:         6041,
		Username:     "root",
		Password:     "taosdata",
		Db:           "",
		MaxIdleConns: 10,
		MaxIdleTime:  0,
		MaxLifeTime:  0,
		MaxOpenConns: 10,
	}
}
