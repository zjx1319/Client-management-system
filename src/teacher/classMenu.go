package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

//显示课堂菜单
func ShowClassMenu() {
	for {
		color.Cyan("┏━━━━━━━━━━━ 课堂菜单━━━━━━━━━━━┓\n")
		color.Cyan("┃         1  进入课堂           ┃\n")
		color.Cyan("┃         2  课堂数据           ┃\n")
		color.Cyan("┃         3  作业数据           ┃\n")
		color.Cyan("┃         4  学生数据           ┃\n")
		color.Cyan("┃         5  退出系统           ┃\n")
		color.Cyan("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
		color.Cyan("请选择（1-5）：\n")

		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			enterClass()
		case 2:
			checkClassData()
		case 3:
			checkWorkData()
		case 4:
			checkUserData()
		case 5:
			os.Exit(0)
		default:
			color.Cyan("输入有误，请重新选择\n")
		}
	}
}

//显示进入课堂后菜单
func EnterClassMenu() {
	for {
		color.Cyan("┏━━━━━━━━━━━ 进入课堂━━━━━━━━━━━┓\n")
		color.Cyan("┃         1  开始上课           ┃\n")
		color.Cyan("┃         2  发送消息           ┃\n")
		color.Cyan("┃         3  课堂数据           ┃\n")
		color.Cyan("┃         4  作业数据           ┃\n")
		color.Cyan("┃         5  学生数据           ┃\n")
		color.Cyan("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
		color.Cyan("请选择（1-5）：\n")

		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			beginClass()
		case 2:
			sendMsg("")
		case 3:
			checkClassData()
		case 4:
			checkWorkData()
		case 5:
			checkUserData()
		default:
			color.Cyan("输入有误，请重新选择\n")
		}
	}
}

//显示开始上课后菜单
func BeginClassMenu() {
	for {
		color.Cyan("┏━━━━━━━━━━━ 正在上课━━━━━━━━━━━┓\n")
		color.Cyan("┃         1  发送消息           ┃\n")
		color.Cyan("┃         2  查看屏幕           ┃\n")
		color.Cyan("┃         3  课堂数据           ┃\n")
		color.Cyan("┃         4  作业数据           ┃\n")
		color.Cyan("┃         5  发布作业           ┃\n")
		color.Cyan("┃         6  学生数据           ┃\n")
		color.Cyan("┃         7  结束下课           ┃\n")
		color.Cyan("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
		color.Cyan("请选择（1-7）：\n")

		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			sendMsg("")
		case 2:
			checkScreen()
		case 3:
			checkClassData()
		case 4:
			checkWorkData()
		case 5:
			addWork()
		case 6:
			checkUserData()
		case 7:
			endClass()
		default:
			color.Cyan("输入有误，请重新选择\n")
		}
	}
}

//显示下课后的菜单
func EndClassMenu() {
	for {
		color.Cyan("┏━━━━━━━━━━━ 课堂结束━━━━━━━━━━━┓\n")
		color.Cyan("┃         1  退出系统           ┃\n")
		color.Cyan("┃         2  课堂数据           ┃\n")
		color.Cyan("┃         3  作业数据           ┃\n")
		color.Cyan("┃         4  学生数据           ┃\n")
		color.Cyan("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
		color.Cyan("请选择（1-4）：\n")
		var key int
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			if len(userMgr) != 0 {
				color.Cyan("当前仍有学生在线，您确定要退出吗\n")
				color.Cyan("输入0确认退出 其他数字返回菜单：\n")
				fmt.Scanf("%d\n", &key)
				if key == 0 {
					os.Exit(0)
				}
			} else {
				os.Exit(0)
			}

		case 2:
			checkClassData()
		case 3:
			checkWorkData()
		case 4:
			checkUserData()
		default:
			color.Cyan("输入有误，请重新选择\n")
		}
	}
}
