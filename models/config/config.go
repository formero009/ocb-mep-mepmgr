package config

type ParameterUpdate struct {
	ParaGroup    string `json:"paraGroup"`										// 参数组名称
	ParaName     string `json:"paraName"`										// 参数组名称
	ParaValue    string `json:"paraValue"`										// 参数值
	ConfigFileId string `json:"configFileId"`									// 对应的配置文件ID
}
