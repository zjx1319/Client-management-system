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
		default:
			fmt.Printf("%s 消息类型无法处理 请联系老师解决\n", msg.Type)
		}
	}
}
