package main

import (
	"fmt"
)

func checkScreen() {
	//输出学生屏幕
	fmt.Println("当前在线学生屏幕信息如下")
	fmt.Println("学号\t姓名\t屏幕未改变时间")
	for userId := range userMgr {
		fmt.Printf("%s\t%s\t%d分钟\n", userMgr[userId].userData.UserId,
			userMgr[userId].userData.UserName, userMgr[userId].screenUnchangeTime)
	}
	//输出完成
}
