package data

import "time"

//课程数据
type Class struct {
	ClassId   int    `json:"classId"`   //课程ID
	ClassName string `json:"className"` //课程名
	ClassNo   int    `json:"classNo"`   //上课总数
	WorkNo    int    `json:"workNo"`    //作业总数
}

//一节课的数据
type ClassData struct {
	ClassNo   int       `json:"classNo"`   //课程numero sign
	BeginTime time.Time `json:"beginTime"` //开始时间
	EndTime   time.Time `json:"endTime"`   //结束时间
}
