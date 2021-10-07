package main

import (
	"bufio"
	"encoding/json"
	"os"
	"src/data"
	"src/tcp"

	"github.com/fatih/color"
)

func sendMsg() {
	var msg data.Message
	var dataByte []byte
	for {
		color.Cyan("请输入你想发送的内容 输入“.exit” 返回主菜单\n")
		color.Cyan("输入“.private”可发送私聊消息\n")
		color.Cyan("请输入内容：\n")
		input := bufio.NewScanner(os.Stdin)
		for {
			input.Scan()
			if input.Text() == ".exit" {
				return
			} else if input.Text() == ".private" {
				var chatPMes data.ChatPMes
				color.Cyan("请输入对方ID(输入teacher可发给老师)\n")
				input.Scan()
				chatPMes.RecieveId = input.Text()
				color.Cyan("请输入发送的内容\n")
				input.Scan()
				chatPMes.Content = input.Text()
				//数据处理
				msg.Type = data.ChatPMesType
				dataByte, _ = json.Marshal(chatPMes)
				msg.Data = string(dataByte)
				dataByte, _ = json.Marshal(msg)
				//发送数据
				tcp.WritePkg(conn, dataByte)

				color.HiMagenta("[Private][%s][to %s]%s:%s\n", user.UserId, chatPMes.RecieveId, user.UserName, chatPMes.Content)
			} else if input.Text() == "" {
				color.Cyan("不能发送空消息哦\n")
			} else {
				//数据处理
				var chatMes data.ChatMes
				msg.Type = data.ChatMesType
				chatMes.Content = input.Text()
				dataByte, _ = json.Marshal(chatMes)
				msg.Data = string(dataByte)
				dataByte, _ = json.Marshal(msg)
				//发送数据
				tcp.WritePkg(conn, dataByte)

				color.HiBlue("[Message][%s]%s:%s\n", user.UserId, user.UserName, chatMes.Content)
			}
		}
	}
}
