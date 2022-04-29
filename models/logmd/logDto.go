package logmd

/**
  "log_time" : 1612687695412,
  "user_name" : "postgres",
  "@timestamp" : "2021-02-07T08:48:15.412Z",
  "database_name" : "ecpserver",
  "application_name" : "Navicat",
  "remote_addr" : "192.168.200.170:7926",
  "level" : "INFO",
  "type" : "database_log",
  "location" : null,
  "query" : null,
  "@version" : "1",
  "message" : "statement: update ko_plugin set ko_route_id = '1' where ko_plugin_id = '3f74f401-1a09-4110-befb-e6f52ae9d3cc'"
*/
type DatabaseLog struct {
	LogTime         int64  `json:"log_time" csv:"-" mapstructure:"log_time"`
	Time            string `json:"-" csv:"log_time" mapstructure:"-"`
	UserName        string `json:"user_name" csv:"user_name" mapstructure:"user_name"`
	DatabaseName    string `json:"database_name" csv:"database_name" mapstructure:"database_name"`
	ApplicationName string `json:"application_name" csv:"application_name" mapstructure:"application_name"`
	RemoteAddr      string `json:"remote_addr" csv:"remote_addr" mapstructure:"remote_addr"`
	Level           string `json:"level" csv:"level" mapstructure:"level"`
	Type            string `json:"type" csv:"type" mapstructure:"type"`
	Location        string `json:"location" csv:"location" mapstructure:"location"`
	Query           string `json:"query" csv:"query" mapstructure:"query"`
	Message         string `json:"message" csv:"message" mapstructure:"message"`
}

type SecurityLog struct {
	LogTime int64  `json:"log_time" csv:"-" mapstructure:"log_time"`
	Time    string `json:"-" csv:"log_time" mapstructure:"-"`
	Account string `json:"account" csv:"account" mapstructure:"account"`
	User    string `json:"user" csv:"user" mapstructure:"user"`
	Level   string `json:"level" csv:"level" mapstructure:'level'`
	OptType string `json:"opt_type" csv:"opt_type" mapstructure:"opt_type"`
	Role    string `json:"role" csv:"role" mapstructure:"role"`
	Ip      string `json:"ip" csv:"ip" mapstructure:"ip"`
}

type SecurityLogResp struct {
	PageNo     int64         `json:"page_no"`
	PageSize   int64         `json:"page_size"`
	TotalPage  int64         `json:"total_page" `
	TotalCount int64         `json:"total_count"`
	List       []SecurityLog `json:"list"`
}

/**
  "log_time" : 1612659600001,
  "level" : "INFO",
  "type" : "system_log",
  "@timestamp" : "2021-02-07T01:00:00.001Z",
  "service_name" : "auth",
  "@version" : "1",
  "msg" : "[pool-1-thread-1] c.c.i.a.s.impl.usermodule.user.UserServiceManager(line:583) - [Unlock]auto unlock user begin."
*/
type SystemLog struct {
	LogTime     int64  `json:"log_time" csv:"-" mapstructure:"log_time"`
	Time        string `json:"-" csv:"log_time" mapstructure:"-"`
	Level       string `json:"level" csv:"level" mapstructure:"level"`
	Type        string `json:"type" csv:"type" mapstructure:"type"`
	ServiceName string `json:"service_name" csv:"service_name" mapstructure:"service_name"`
	Message     string `json:"msg" csv:"msg" mapstructure:"msg"`
}

type NeOperationLog struct {
	LogTime    int64  `json:"log_time" csv:"-" mapstructure:"log_time"`
	Time       string `json:"-" csv:"log_time" mapstructure:"-"`
	Account    string `json:"account" csv:"account" mapstructure:"account"`
	Name       string `json:"name" csv:"name" mapstructure:"name"`
	Role       string `json:"role" csv:"role" mapstructure:"role"`
	Result     string `json:"result" csv:"result" mapstructure:"result"`
	Ip         string `json:"ip" csv:"ip" mapstructure:"ip"`
	OptContent string `json:"opt_content" csv:"opt_content" mapstructure:"opt_content"`
	NeName     string `json:"ne_name" csv:"ne_name" mapstructure:"ne_name"`
}

type NeOperationLogResp struct {
	PageNo     int64            `json:"page_no"`
	PageSize   int64            `json:"page_size"`
	TotalPage  int64            `json:"total_page"`
	TotalCount int64            `json:"total_count"`
	List       []NeOperationLog `json:"list"`
}

type UploadReq struct {
	Address   string        `json:"address" required:"true" description:"ftp server address"`
	Username  string        `json:"username" required:"true" description:"ftp server username"`
	Password  string        `json:"password" required:"true" description:"ftp server password"`
	Type      string        `json:"type" required:"true" description:"upload type:all/selected/self"`
	Data      []interface{} `json:"data" required:"false" description:"if the type is selected, then it represent the selected data"`
	StartTime int64         `json:"start_time" required:"false" description:"if the type is self, then it represent the start time of query"`
	EndTime   int64         `json:"end_time" required:"false" description:"if the type is self, then it represent the end time of query"`
}

type ExportReq struct {
	Data      []interface{} `json:"data" required:"false" description:"selected range data"`
	Format    string        `json:"format" required:"true" description:"exported file format:txt/csv/xlsx"`
	Range     string        `json:"range" required:"true" description:"export range:all/selected/self"`
	StartTime int64         `json:"start_time" required:"false" description:"if the range is self, then it represent the start time of query"`
	EndTime   int64         `json:"end_time" required:"false" description:"if the range is self, then it represent the end time of query"`
}

type ExportDatabaseLogResp struct {
	Count int64          `json:"count"`
	Data  []*DatabaseLog `json:"data"  description:"database log list"`
}
type ExportSystemLogResp struct {
	Count int64        `json:"count"`
	Data  []*SystemLog `json:"data" description:"software system log list"`
}
type ExportSecurityLogResp struct {
	Count int64          `json:"count"`
	Data  []*SecurityLog `json:"data"  description:"security log list"`
}
type ExportNeOperationLogResp struct {
	Count int64             `json:"count"`
	Data  []*NeOperationLog `json:"data"  description:"network operation log list"`
}
