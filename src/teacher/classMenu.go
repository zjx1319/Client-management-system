package main

import (
	"fmt"
	"os"
)

func ShowClassMenu() {
	for {
		fmt.Println("------------------")
		fmt.Println("1. 进入课堂")
		fmt.Println("2. 课堂数据")
		fmt.Println("3. 作业数据")
		fmt.Println("4. 退出系统")
		fmt.Println("请输入：")
		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("进入课堂")
			enterClass()
		case 2:
			fmt.Println("课堂数据")
			checkClassData()
		case 3:
			fmt.Println("作业数据")
		case 4:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("输入有误，请重新选择")
		}
	}
}

func EnterClassMenu() {
	for {
		fmt.Println("------------------")
		fmt.Println("1. 开始上课")
		fmt.Println("2. 发送消息")
		fmt.Println("3. 课堂数据")
		fmt.Println("4. 作业数据")
		fmt.Println("请输入：")
		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("开始上课")
			sendMsg("同学们开始上课啦~")
			beginClass()
		case 2:
			fmt.Println("发送消息")
			sendMsg("")
		case 3:
			fmt.Println("课堂数据")
			checkClassData()
		case 4:
			fmt.Println("作业数据")
		default:
			fmt.Println("输入有误，请重新选择")
		}
	}
}

func BeginClassMenu() {
	for {
		fmt.Println("------------------")
		fmt.Println("1. 发送消息")
		fmt.Println("2. 课堂数据")
		fmt.Println("3. 作业数据")
		fmt.Println("4. 发布作业")
		fmt.Println("5. 结束下课")
		fmt.Println("请输入：")
		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("发送消息")
			sendMsg("")
		case 2:
			fmt.Println("课堂数据")
			checkClassData()
		case 3:
			fmt.Println("作业数据")
		case 4:
			fmt.Println("发布作业")
		case 5:
			fmt.Println("结束下课")
			sendMsg("下课啦 请同学们尽快退出登录~")
			endClass()
		default:
			fmt.Println("输入有误，请重新选择")
		}
	}
}

func EndClassMenu() {
	for {
		fmt.Println("------------------")
		fmt.Println("课程已结束 请提醒同学们尽快退出登录")
		fmt.Println("建议等待同学们全部退出后再关闭教师端")
		fmt.Println("1. 退出系统")
		fmt.Println("请输入：")
		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("输入有误，请重新选择")
		}
	}
}
