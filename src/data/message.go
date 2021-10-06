package data

import "time"

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

	//消息
	ChatMesType     = "ChatMes"
	ChatResMesType  = "ChatResMes"
	ChatPMesType    = "ChatPMes"
	ChatPResMesType = "ChatPResMes"

	//提问
	QueMesType = "QueMes"

	//作业
	WorkAllMesType    = "WorkAllMes" //获取作业数量
	WorkAllResMesType = "WorkAllResMes"
	WorkMesType       = "WorkMes" //获取作业信息
	WorkResMesType    = "WorkResMes"
	WorkSubMesType    = "WorkSubMes" //作业提交信息
	WorkSubResMesType = "WorkSubResMes"
	Work_Objective    = 1 //客观题
	Work_Subjective   = 2 //主观题
	Work_Files        = 3 //文件题

	//屏幕相关
	ScreenReportType     = "ScreenReport"  //屏幕内容报告
	ScreenShotGetType    = "ScreenShotGet" //获取屏幕截图
	ScreenShotResType    = "ScreenShotRes" //返回屏幕截图
	ScreenVideoStartType = "ScreenVideoStart"
	ScreenVideoStopType  = "ScreenVideoStop"
	ScreenVideoResType   = "ScreenVideoRes"
)

//消息结构体，含有两个部分,消息的类别和消息的内容
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type FileMes struct {
	FileName string `json:"fileName"`
	Data     string `json:"data"` //使用base64编码
}

//消息种类：ScreenVideoRes 包含屏幕图片
type ScreenVideoRes struct {
	Img string `json:"img"` //使用base64编码
}

//消息种类：ScreenReport 包含屏幕内容未改变时长（单位分钟）
type ScreenReport struct {
	UnchangeTime int `json:"unchangeTime"`
}

//消息种类：ScreenShotGet
type ScreenShotGet struct {
}

//消息种类：ScreenShotRes 包含屏幕截图
type ScreenShotRes struct {
	Img string `json:"img"` //使用base64编码
}

//消息种类：WorkAllMes
type WorkAllMes struct {
}

//消息种类：WorkAllResMes  包含作业总数
type WorkAllResMes struct {
	Num int `json:"num"`
}

//消息种类：WorkMes 作业信息，包含作业ID
type WorkMes struct {
	Id int `json:"id"`
}

//消息种类：WorkResMes 作业信息，包含作业ID、答案、提交时间、分数
type WorkResMes struct {
	Id         int       `json:"id"`
	Type       int       `json:"type"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"` //如果答案为空则未提交
	SubmitTime time.Time `json:"submitTime"`
	DeadLine   time.Time `json:"deadLine"`
	Score      int       `json:"score"`
	FullScore  int       `json:"fullScore"`
}

//消息种类：WorkSubMes 作业提交信息，包含作业ID、答案
type WorkSubMes struct {
	Id     int    `json:"id"`
	Answer string `json:"answer"`
}

//消息种类：WorkSubResMes 作业提交信息，包含作业ID、分数（客观题）
type WorkSubResMes struct {
	Id    int `json:"id"`
	Score int `json:"score"`
}

//消息种类：QueMes 提问信息，包含提问内容
type QueMes struct {
	Content string `json:"content"`
}

//消息种类：ChatMes 聊天消息，包含消息内容
type ChatMes struct {
	Content string `json:"content"`
}

//消息种类：ChatResMes 聊天消息，包含发送者id 姓名 消息内容
type ChatResMes struct {
	SendUserId   string `json:"sendUserId"`
	SendUserName string `json:"sendUserName"`
	Content      string `json:"content"`
}

//消息种类：ChatPMes 私聊消息，包含接收id，消息内容
type ChatPMes struct {
	RecieveId string `json:"recieveId"`
	Content   string `json:"content"`
}

//消息种类：ChatPResMes 聊天消息，包含发送者id 姓名 消息内容
type ChatPResMes struct {
	SendUserId   string `json:"sendUserId"`
	SendUserName string `json:"sendUserName"`
	Content      string `json:"content"`
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
