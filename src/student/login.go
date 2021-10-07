package main

import (
	"encoding/json"
	"net"
	"src/config"
	"src/data"
	"src/tcp"

	"github.com/fatih/color"
)

func login() (err error) {
	//建立连接
	conn, err = net.Dial("tcp", config.Server)
	if err != nil {
		return
	}
	//处理信息
	var msg data.Message
	msg.Type = data.LoginMesType
	var loginMes data.LoginMes
	loginMes.UserId = user.UserId
	loginMes.UserPwd = user.UserPwd
	loginMes.Seat = config.Seat

	//使用json序列化
	dataByte, err := json.Marshal(loginMes)
	if err != nil {
		return
	}
	msg.Data = string(dataByte)
	dataByte, err = json.Marshal(msg)
	if err != nil {
		return
	}
	//发送消息给服务器
	err = tcp.WritePkg(conn, []byte(dataByte))
	if err != nil {
		return
	}
	//读服务器返回的消息
	msg, err = tcp.ReadPkg(conn)

	//json反序列化
	var loginResMes data.LoginResMes
	json.Unmarshal([]byte(msg.Data), &loginResMes)
	if loginResMes.Result == data.Login_IDNotFound {
		color.Red("用户不存在!\n")
		return
	} else if loginResMes.Result == data.Login_PwdError {
		color.Red("密码错误！\n")
		return
	} else if loginResMes.Result == data.Login_Success {
		user.UserName = loginResMes.Username
		color.Green("登录成功,%s,欢迎您\n", user.UserName)
		color.Cyan("当前课程为：%s \n", loginResMes.ClassName)
		//启动消息接收线程
		go process()
		//启动屏幕内容检测线程
		go screenCheckProcess()
		//显示菜单
		ShowClassMenu()
	}
	return
}
