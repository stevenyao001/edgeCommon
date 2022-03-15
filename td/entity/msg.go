package entity

import (
	"github.com/stevenyao001/edgeCommon/td"
)

type UserTd struct {
	*td.Engine
	MegnetStatus bool `json:"megnet_status"`
	Ia           int  `json:"ia"`
	Ep           int  `json:"ep"`
}

type UserTdGK struct {
	*td.Engine
	MegnetStatus bool `json:"megnet_status"`
	Num          int  `json:"num"`
	Ia           int  `json:"ia"`
	Status       int  `json:"status"`
	Ep           int  `json:"ep"`
	InstantEp    int  `json:"instant_ep"`
}
