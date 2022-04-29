package common

const (
	ErrBasicPrefix = "opensimga.mepm.mepmgr."

	ErrSuccess      = "opensigma.common.common.Success"
	ErrParaJson     = ErrBasicPrefix + "InvalidJson"
	ErrParaInvalid  = ErrBasicPrefix + "InvalidParam"
	ErrParaDupli    = ErrBasicPrefix + "DuplicateParam"
	ErrAlreadyExist = ErrBasicPrefix + "AlreadyExist"
	ErrNotFound     = ErrBasicPrefix + "NotFound"
	ErrDatabase     = ErrBasicPrefix + "DatabaseError"
	ErrNetwork      = ErrBasicPrefix + "NetworkError"
	ErrMepStatus    = ErrBasicPrefix + "MepStatus"
	ErrMepResp      = ErrBasicPrefix + "MepError"
	ErrCannotDel    = ErrBasicPrefix + "CannotDelete"
	ErrInternal     = ErrBasicPrefix + "InternalError"
	ErrAuthFailed   = ErrBasicPrefix + "AuthFailed"
	ErrAuthResp     = ErrBasicPrefix + "AuthResp"


	ErrUnKnown = ErrBasicPrefix + "UnKnown"

)

const (
	MsgLackRequiredPara = "缺少必填参数，请检查"
	MsgInvalidPara      = "取值非法"
	MsgDuplicatePara    = "取值重复"
	MsgInvalidFormat    = "取值格式错误"
	MsgAlreadyExist     = "已存在"
	MsgNotFound         = "不存在"
	MsgCannotDel        = "不能删除"
)

var code2Msg = map[string]string{
	ErrSuccess:      "请求成功",
	ErrParaJson:     "参数格式错误",
	ErrParaInvalid:  "参数非法",
	ErrAlreadyExist: "已存在",
	ErrNotFound:     "不存在",
	ErrParaDupli:    "参数取值重复",
	ErrDatabase:     "数据库错误",
	ErrNetwork:      "网络异常",
	ErrMepResp:      "mep错误",
	ErrMepStatus:    "mep状态异常",
	ErrCannotDel:    "不能删除",
	ErrInternal:     "内部错误",

	ErrUnKnown: "未定义错误",
}

type ErrMsg struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrMsg) Error() string {
	return e.Message
}

func NewError(code string, params ...string) error {
	message := ""
	for _, m := range params {
		if m != "" {
			message = m
			break
		}
	}
	if message == "" {
		ok := false
		if message, ok = code2Msg[code]; !ok {
			message = ErrUnKnown
		}
	}
	return ErrMsg{Code: code, Message: message}
}
