package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"src/data"
	"src/tcp"
)

func sendMsg() {
	var msg data.Message
	var dataByte []byte
	for {
		fmt.Println("------------------")
		fmt.Println("请输入你想发送的内容 输入“.exit” 返回主菜单")
		fmt.Println("输入“.private”可发送私聊消息 ")
		fmt.Println("请输入内容：")
		input := bufio.NewScanner(os.Stdin)
		for {
			input.Scan()
			if input.Text() == ".exit" {
				return
			} else if input.Text() == ".private" {
				var chatPMes data.ChatPMes
				fmt.Println("请输入对方ID(输入teacher可发给老师)")
				input.Scan()
				chatPMes.RecieveId = input.Text()
				fmt.Println("请输入发送的内容")
				input.Scan()
				chatPMes.Content = input.Text()
				//数据处理
				msg.Type = data.ChatPMesType
				dataByte, _ = json.Marshal(chatPMes)
				msg.Data = string(dataByte)
				dataByte, _ = json.Marshal(msg)
				//发送数据
				tcp.WritePkg(conn, dataByte)

				fmt.Printf("[P][%s][to %s]%s:%s\n", user.UserId, chatPMes.RecieveId, user.UserName, chatPMes.Content)
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

				fmt.Printf("[%s]%s:%s\n", user.UserId, user.UserName, chatMes.Content)
			}
		}
	}
}
