package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"src/data"
	"src/tcp"
)

//处理和客户端的通讯
func process() (err error) {
	for {
		var msg data.Message
		msg, err = tcp.ReadPkg(conn)
		if err != nil {
			// 通常遇到的错误是连接中断或被关闭，用io.EOF表示
			if err == io.EOF {
				if user.UserName != "" {
					fmt.Println("异常掉线！请重新登录！")
				}
			} else {
				fmt.Println(err)
			}
			return
		}
		switch msg.Type {
		case data.LogoutResMesType:
			//处理注销
			//使用json序列化
			var logoutResMes data.LogoutResMes
			json.Unmarshal([]byte(msg.Data), &logoutResMes)
			if logoutResMes.Result == data.Logout_Success {
				conn.Close()
				fmt.Println("已退出课堂")
				os.Exit(0)
			} else {
				fmt.Printf("注销失败 返回信息%s\n", logoutResMes.Result)
			}
		case data.ChatResMesType:
			//接收到消息
			//使用json序列化
			var chatResMes data.ChatResMes
			json.Unmarshal([]byte(msg.Data), &chatResMes)
			fmt.Printf("[M][%s]%s:%s\n", chatResMes.SendUserId, chatResMes.SendUserName, chatResMes.Content)
		case data.ChatPResMesType:
			//接收到私聊消息
			//使用json序列化
			var chatPResMes data.ChatPResMes
			json.Unmarshal([]byte(msg.Data), &chatPResMes)
			fmt.Printf("[P][%s]%s:%s\n", chatPResMes.SendUserId, chatPResMes.SendUserName, chatPResMes.Content)
		case data.WorkAllResMesType:
			var workAllResMes data.WorkAllResMes
			json.Unmarshal([]byte(msg.Data), &workAllResMes)
			go getWorkData(workAllResMes.Num)
		case data.WorkResMesType:
			var workResMes data.WorkResMes
			json.Unmarshal([]byte(msg.Data), &workResMes)
			go workDataRes(workResMes)
		case data.WorkSubResMesType:
			var workSubResMes data.WorkSubResMes
			json.Unmarshal([]byte(msg.Data), &workSubResMes)
			subWorkRes(workSubResMes)
		case data.ScreenShotGetType:
			sendScreenShot()
		case data.ScreenVideoStartType:
			VideoFlag = true
			go sendScreenVideo()
		case data.ScreenVideoStopType:
			VideoFlag = false
		default:
			fmt.Printf("%s 消息类型无法处理\n", msg.Type)
		}
	}
}
