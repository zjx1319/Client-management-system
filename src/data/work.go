package data

import "time"

//作业数据
type WorkData struct {
	Id        int       `json:"id"`
	Type      int       `json:"type"`
	Question  string    `json:"question"`
	StdAnswer string    `json:"stdAnswer"` //标准答案
	FullScore int       `json:"fullScore"`
	DeadLine  time.Time `json:"deadLine"`
}
