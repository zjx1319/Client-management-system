package main

import (
	"fmt"
	"os"
	"src/config"

	"github.com/fatih/color"
	"github.com/garyburd/redigo/redis"
)

var RconnPool *redis.Pool

//初始化redis连接池
func init() {
	RconnPool = &redis.Pool{
		MaxIdle:     config.RedisMaxIdle,
		MaxActive:   config.RedisMaxActive,
		IdleTimeout: config.RedisIdleTimeout,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", config.RedisIp)
			if err != nil {
				return
			}
			if class.ClassId != 0 {
				conn.Do("select", class.ClassId) //切换到课程数据库
			}
			return
		},
	}
}

func main() {
	var key int

	for {
		color.Cyan("┏ 欢迎使用电子教室管理系统教师端┓\n")
		color.Cyan("┃         1  选择课程           ┃\n")
		color.Cyan("┃         2  添加课程           ┃\n")
		color.Cyan("┃         3  退出系统           ┃\n")
		color.Cyan("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
		color.Cyan("请选择（1-3）：\n")

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			selectClass()
		case 2:
			addClass()
		case 3:
			os.Exit(0)
		default:
			color.Cyan("输入有误，请重新选择\n")
		}
	}
}

func errorReport(err error) {
	color.Red("[Error]运行过程中发生错误：%s\n", err.Error())
}
