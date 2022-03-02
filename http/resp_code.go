package http

type RespCode int

const (
	RespCodeSuccess    RespCode = 200
	RespCodeUnknownErr RespCode = 10000
	RespCodeParamErr   RespCode = 10001

	RespCodeUserNotExists RespCode = 10100
)

var RespMsg map[RespCode]string

func init() {
	RespMsg = map[RespCode]string{
		RespCodeSuccess:       "success",
		RespCodeUnknownErr:    "unknown err",
		RespCodeParamErr:      "param err",
		RespCodeUserNotExists: "user not exists",
	}
}
