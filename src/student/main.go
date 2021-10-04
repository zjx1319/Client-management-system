package main

import (
	"fmt"
	"net"
	"os"
	"src/config"
	"src/data"
)

var user data.User
var conn net.Conn

func main() {
	//测试用
	com := os.Args
	if len(com) == 4 {
		user.UserId = com[1]
		user.UserPwd = com[2]
		config.Seat = com[3]
		login()
	}

	for {
		fmt.Println("--------欢迎使用电子教室管理系统学生端-------")
		fmt.Print("请输入学号：")
		fmt.Scanf("%s\n", &user.UserId)
		fmt.Print("请输入密码：")
		fmt.Scanf("%s\n", &user.UserPwd)

		//登录
		login()
	}
}
