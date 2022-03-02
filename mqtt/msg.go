package mqtt

import "sync"

type Msg struct {
	TraceId  string                 `json:"trace_id"`
	MsgId    int64                  `json:"msg_id"`
	DeviceId string                 `json:"device_id"`
	Version  string                 `json:"version"`
	Source   Source                 `json:"source"`
	Mold     CmdMold                `json:"mold"`
	Cmd      Cmd                    `json:"cmd"`
	Content  map[string]interface{} `json:"content"`
}

var MsgPool = sync.Pool{
	New: func() interface{} {
		return Msg{
			TraceId:  "",
			MsgId:    0,
			DeviceId: "",
			Version:  "",
			Source:   0,
			Mold:     0,
			Cmd:      0,
			Content:  nil,
		}
	},
}

//消息来源
type Source int32

const (
	SourceCollect   Source = 100000 //采集器-100000
	SourceEdge      Source = 200000 //边缘计算器-200000
	SourceEdgeAdmin Source = 300000 //边缘计算管理-300000
	SourceIot       Source = 400000 //云端-400000
)

//指令类型
type CmdMold int32

const (
	DeviceControl CmdMold = 1000 //设备控制-10000
	DataHandle    CmdMold = 2000 //数据处理-20000
)

//消息指令
type Cmd int32

const (
	CollectDeviceRegister Cmd = 101001 //采集器-设备-注册（上报）
	CollectDeviceDel      Cmd = 101003 //采集器-设备-销毁（上报）
	CollectDataReport     Cmd = 102001 //采集器-数据-上报（上报）

	EdgeDeviceRegister Cmd = 201001 //边缘计算器-设备-注册（上报）
	EdgeDeviceUpgrade  Cmd = 201002 //边缘计算器-设备-升级（下发）
	EdgeDataReport     Cmd = 202001 //边缘计算器-数据-上报（上报）

	EdgeAdminDeviceRegister Cmd = 301001 //边缘计算管理器-设备-注册（上报）
	EdgeAdminDeviceUpgrade  Cmd = 301002 //边缘计算管理器-设备-升级（下发）
	EdgeAdminDataReport     Cmd = 302001 //边缘计算管理器-数据-上报（上报）

	IotDeviceRegister Cmd = 401001 //云端-设备-注册（下发）
	IotDeviceUpgrade  Cmd = 401002 //云端-设备-升级（下发）
	IotData           Cmd = 402001 //云端-数据-下发
)
