package data

import "time"

const (
	//课堂状态
	ClassStatus_Default    = 0 //默认
	ClassStatus_Normal     = 1 //正常
	ClassStatus_LeaveEarly = 2 //早退
	ClassStatus_JoinLate   = 3 //迟到
	ClassStatus_LEJL       = 4 //迟到早退
	ClassStatus_Absent     = 5 //缺勤
)

//用户数据
type User struct {
	UserId   string `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

//用户课堂数据
type UserClassData struct {
	UserId      string    `json:"userId"`
	ClassStatus int       `json:"classStatus"`
	JoinTime    time.Time `json:"joinTime"`
	LeaveTime   time.Time `json:"leaveTime"`
	Seat        string    `json:"seat"`
}
