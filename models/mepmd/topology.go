package mepmd

type TopologyMepGroupRe struct {
	MepId        string
	MepName      string
	Id           int64
	MepGroupName string
}

type TopologyMep struct {
	MepId   string `json:"mepId"`
	MepName string `json:"mepName"`
}

type RespTopology struct {
	MepGroupId   int64         `json:"mepGroupId"`
	MepGroupName string        `json:"mepGroupName"`
	Meps         []TopologyMep `json:"meps"`
}

type RespTopologyByCity struct {
	City string        `json:"city"`
	Meps         []TopologyMep `json:"meps"`
}

type RespTopologyByProvincial struct {
	Province   string     `json:"province"`
	Cites        []RespTopologyByCity `json:"cites"`
}

type AlarmRec struct {
	NetworkName string `json:"networkName"`
	AlarmResp
}

type AlarmResp struct {
	AlarmNumber    int64 `json:"alarmNumber"`
	CriticalNumber int64 `json:"criticalNumber"`
	MajorNumber    int64 `json:"majorNumber"`
	MinorNumber    int64 `json:"minorNumber"`
	WarningNumber  int64 `json:"warningNumber"`
}

type AlarmStatistics struct {
	Code    string        `json:"code"`
	Message string     `json:"message"`
	Data    []AlarmRec `json:"data"`
}

type RespMepDetail struct {
	MepId      string    `json:"mepId"`
	MepName    string    `json:"mepName"`
	EndPoint   string    `json:"endPoint"`
	User       string    `json:"user"`
	Type       string    `json:"type"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	UserTag    string    `json:"userTag"`
	Longitude  string    `json:"longitude"`
	Latitude   string    `json:"latitude"`
	Contractor string    `json:"contractor"`
	RunStatus  string    `json:"runStatus"`
	MgrStatus  string    `json:"mgrStatus"`
	SwVersion  string    `json:"swVersion"`
	AlarmMsg   AlarmResp `json:"alarm"`
}

type TopologyReq struct {
	CALLERId   string `json:"CALLERId"`
	MEPMID     string `json:"MEPMID"`
	MepOrGroup string `json:"MepOrGroup"`
	MepId      string `json:"mepId"`
	MepGroupId string `json:"mepGroupId"`
	Province string `json:"province"`
	City       string `json:"city"`
}
