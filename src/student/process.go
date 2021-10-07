package main

import (
	"encoding/json"
	"os"
	"src/data"
	"src/tcp"

	"github.com/fatih/color"
)

//处理和客户端的通讯
func process() (err error) {
	for {
		var msg data.Message
		msg, err = tcp.ReadPkg(conn)
		if err != nil {
			color.Red("连接中断！请重新登录！\n")
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
				color.Cyan("已退出课堂")
				os.Exit(0)
			} else {
				color.Red("注销失败 返回信息%s\n", logoutResMes.Result)
			}
		case data.ChatResMesType:
			//接收到消息
			//使用json序列化
			var chatResMes data.ChatResMes
			json.Unmarshal([]byte(msg.Data), &chatResMes)
			color.HiBlue("[Message][%s]%s:%s\n", chatResMes.SendUserId, chatResMes.SendUserName, chatResMes.Content)
		case data.ChatPResMesType:
			//接收到私聊消息
			//使用json序列化
			var chatPResMes data.ChatPResMes
			json.Unmarshal([]byte(msg.Data), &chatPResMes)
			color.HiMagenta("[Private][%s]%s:%s\n", chatPResMes.SendUserId, chatPResMes.SendUserName, chatPResMes.Content)
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
		case data.BlockListProcessType:
			var BlockListProcess data.BlockListProcess
			json.Unmarshal([]byte(msg.Data), &BlockListProcess)
			go blockListCheckProcess(BlockListProcess.List)
		case data.BlockListWebType:
			var BlockListWeb data.BlockListWeb
			json.Unmarshal([]byte(msg.Data), &BlockListWeb)
			go blockListWeb(BlockListWeb.List)
		default:
			color.Red("%s 消息类型无法处理\n", msg.Type)
		}
	}
}
