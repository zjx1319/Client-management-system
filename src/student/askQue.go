package main

import (
	"bufio"
	"encoding/json"
	"os"
	"src/data"
	"src/tcp"

	"github.com/fatih/color"
)

func askQue() {
	var msg data.Message
	var dataByte []byte
	color.Cyan("请输入你遇到的问题:\n")
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

		color.HiMagenta("[Question]你发起了一个提问:%s\n", queMes.Content)
	}
}
