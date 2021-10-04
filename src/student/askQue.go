package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"src/data"
	"src/tcp"
)

func askQue() {
	var msg data.Message
	var dataByte []byte
	fmt.Println("------------------")
	fmt.Println("请输入你遇到的问题，不想输入就直接回车吧~")
	fmt.Println("请输入内容：")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if input.Text() == ".exit" {
		return
	} else {
		//数据处理
		var queMes data.QueMes
		msg.Type = data.QueMesType
		queMes.Content = input.Text()
		dataByte, _ = json.Marshal(queMes)
		msg.Data = string(dataByte)
		dataByte, _ = json.Marshal(msg)
		//发送数据
		tcp.WritePkg(conn, dataByte)

		fmt.Printf("[Q]你发起了一个提问:%s\n", queMes.Content)
	}
}
