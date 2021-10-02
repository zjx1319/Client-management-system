package data

const (
	//登录
	LoginMesType     = "LoginMes"    //发送
	LoginResMesType  = "LoginResMes" //返回
	Login_IDNotFound = "IDNotFound"  //无用户
	Login_PwdError   = "PwdError"    //密码错误
	Login_Success    = "Success"     //成功

	//注销
	LogoutMesType    = "LogoutMes"    //发送
	LogoutResMesType = "LogoutResMes" //返回
	Logout_Success   = "Success"      //成功
)

//消息结构体，含有两个部分,消息的类别和消息的内容
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

//消息种类: LoginMes 登录消息, 包含id 密码和机位
type LoginMes struct {
	UserId  string `json:"userId"`
	UserPwd string `json:"userPwd"`
	Seat    string `json:"seat"`
}

//消息种类: LoginResMes 登录消息, 包含结果 课程名 用户名
type LoginResMes struct {
	Result    string `json:"Result"`
	ClassName string `json:"classname"`
	Username  string `json:"username"`
}

//消息种类: Loginout 注销消息, 包含id
type LogoutMes struct {
	UserId string `json:"userId"`
}

//消息种类: LogoutResMes 注销消息, 包含结果
type LogoutResMes struct {
	Result string `json:"Result"`
}
