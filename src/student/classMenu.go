package main

import (
	"fmt"
)

func ShowClassMenu() {
	for {
		fmt.Println("------------------")
		fmt.Println("1. 发送消息")
		fmt.Println("2. 举手提问")
		fmt.Println("3. 作业中心")
		fmt.Println("4. 个人中心")
		fmt.Println("5. 退出登录")
		fmt.Println("请输入：")
		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("1. 发送消息")
			sendMsg()
		case 2:
			fmt.Println("2. 举手提问")
			askQue()
		case 3:
			fmt.Println("3. 作业中心")
		case 4:
			fmt.Println("4. 个人中心")
		case 5:
			logout()
		default:
			fmt.Println("输入有误，请重新选择")
		}
	}
}
