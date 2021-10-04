package main

import (
	"fmt"
	"os"
	"src/config"
	"src/data"

	"github.com/garyburd/redigo/redis"
)

var Rconn redis.Conn
var class data.Class
var classData data.ClassData

func init() {
	//当服务器启动时，连接redis
	Rconn, _ = redis.Dial("tcp", config.RedisIp)

	//defer conn.Close()
}

func main() {
	var key int

	for {
		fmt.Println("--------欢迎使用电子教室管理系统教师端-------")
		fmt.Println("                1  选择课程                ")
		fmt.Println("                2  添加课程                ")
		fmt.Println("                3  退出系统                ")
		fmt.Println("请选择（1-3）：")

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("选择课程")
			selectClass()
		case 2:
			fmt.Println("添加课程")
			addClass()
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("输入有误，请重新选择")
		}
	}
}
