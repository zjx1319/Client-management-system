package main

import (
	"fmt"

	"github.com/fatih/color"
)

func ShowClassMenu() {
	for {
		color.Cyan("┏━━━━━━━━━━━ 学生菜单━━━━━━━━━━━┓\n")
		color.Cyan("┃         1  发送消息           ┃\n")
		color.Cyan("┃         2  举手提问           ┃\n")
		color.Cyan("┃         3  作业中心           ┃\n")
		color.Cyan("┃         4  退出登录           ┃\n")
		color.Cyan("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
		color.Cyan("请选择（1-4）：\n")

		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			sendMsg()
		case 2:
			askQue()
		case 3:
			getWorkAll()
		case 4:
			logout()
		default:
			color.Cyan("输入有误，请重新选择\n")
		}
	}
}
