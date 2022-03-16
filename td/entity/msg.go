package entity

type UserTd struct {
	MegnetStatus bool  `json:"megnet_status"`
	Ia           int   `json:"ia"`
	Ep           int64 `json:"ep"`
}

type UserTdGK struct {
	MegnetStatus bool `json:"megnet_status"`
	Num          int  `json:"num"`
	Ia           int  `json:"ia"`
	Status       int  `json:"status"`
	Ep           int  `json:"ep"`
	InstantEp    int  `json:"instant_ep"`
}
