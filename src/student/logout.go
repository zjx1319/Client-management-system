package main

import (
	"encoding/json"
	"fmt"
	"src/data"
	"src/tcp"
)

func logout() (err error) {

	//处理信息
	var msg data.Message
	msg.Type = data.LogoutMesType
	var LogoutMes data.LogoutMes
	LogoutMes.UserId = user.UserId

	//使用json序列化
	dataByte, err := json.Marshal(LogoutMes)
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
	fmt.Println("已发送注销请求 成功前请不要关闭本程序")
	return
}
