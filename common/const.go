package common

const Unknown = "unknown"
const UserInfo = "userInfo"


const (
	AddMep          = "新增mep"
	DeleteMep       = "删除mep"
	UpdateMep       = "更新mep"
	GetMep          = "查询mep"
	ListMep         = "查询mep列表"
	GetMepTopo      = "查询mep拓扑"
	UpdateMepStatus = "更新mep状态"
)

const (
	GetAppList     = "查询APP调用统计"
	GetAppChart    = "查询APP调用统计图"
	GetAppLog      = "查询调用日志"
	GetServiceList = "查询服务调用统计"
)

const (
	MepRootPath      = "/mep/mec_service_mgmt"
	MepVersion1      = "/v1"
	MepInfoPath      = MepRootPath + MepVersion1 + "/info"
	MepAuthPath      = MepRootPath + MepVersion1 + "/authentication"
	MepKeepAlivePath = MepAuthPath + "/alive"

	AuthMepAuthorityPath = "/auth/user/ne_authority"
)

const (
	AuthHeader = "Authorization"
)

const (
	AlarmStatisticsPath = "/mepm/alarm/v1/active/statistics"
)
