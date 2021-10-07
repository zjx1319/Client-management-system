package main

import (
	"fmt"
	"net"
	"os"
	"src/config"
	"src/data"

	"github.com/fatih/color"
)

var user data.User
var conn net.Conn

func main() {
	//读取命令行参数
	com := os.Args
	if len(com) == 4 {
		user.UserId = com[1]
		user.UserPwd = com[2]
		config.Seat = com[3]
		login()
	}
	if len(com) == 2 {
		config.Seat = com[1]
	}

	for {
		color.Cyan("欢迎使用电子教室管理系统学生端")
		color.Cyan("请输入学号：")
		fmt.Scanf("%s\n", &user.UserId)
		color.Cyan("请输入密码：")
		fmt.Scanf("%s\n", &user.UserPwd)

		//登录
		login()
	}
}
