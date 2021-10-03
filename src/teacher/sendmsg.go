package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"src/data"
	"src/tcp"
)

var teacher data.User

func init() {
	teacher.UserId = "teacher"
	teacher.UserName = "teacher"
}
func sendMsg(s string) { //本函数参数为空文本进入菜单 非空可直接发送全体消息
	if s == "" {
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
					fmt.Println("请输入对方ID")
					input.Scan()
					chatPMes.RecieveId = input.Text()
					fmt.Println("请输入发送的内容")
					input.Scan()
					chatPMes.Content = input.Text()
					sendPResMsg(teacher, chatPMes)
				} else {
					sendMsg(input.Text())
				}
			}
		}
	} else {
		//数据处理
		var chatMes data.ChatMes
		chatMes.Content = s
		//发送数据
		sendResMsg(teacher, chatMes)
	}
}

func sendResMsg(user data.User, chatMes data.ChatMes) {
	var chatResMes data.ChatResMes
	var msg data.Message
	chatResMes.SendUserId = user.UserId
	chatResMes.SendUserName = user.UserName
	chatResMes.Content = chatMes.Content
	msg.Type = data.ChatResMesType
	dataByte, _ := json.Marshal(chatResMes)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)

	for _, up := range userMgr {
		if up.userData.UserId != user.UserId {
			tcp.WritePkg(up.conn, dataByte)
		}
	}

	fmt.Printf("[M][%s]%s:%s\n", chatResMes.SendUserId, chatResMes.SendUserName, chatResMes.Content)
}

func sendPResMsg(user data.User, chatPMes data.ChatPMes) {
	//处理发给老师的消息
	if chatPMes.RecieveId == "teacher" {
		fmt.Printf("[P][%s]%s:%s\n", user.UserId, user.UserName, chatPMes.Content)
		return
	}

	//处理一般私聊消息
	var chatPResMes data.ChatPResMes
	var msg data.Message
	chatPResMes.SendUserId = user.UserId
	chatPResMes.SendUserName = user.UserName
	chatPResMes.Content = chatPMes.Content
	msg.Type = data.ChatPResMesType
	dataByte, _ := json.Marshal(chatPResMes)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)

	for _, up := range userMgr {
		if up.userData.UserId == chatPMes.RecieveId {
			tcp.WritePkg(up.conn, dataByte)
		}
	}
}
