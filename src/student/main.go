package main

import (
	"fmt"
	"net"
	"src/data"
)

var user data.User
var conn net.Conn

func main() {
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
