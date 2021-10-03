package main

import (
	"encoding/json"
	"src/data"
	"src/tcp"
	"time"
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
	time.Sleep(1 * time.Minute)
	return
}
