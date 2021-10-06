package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"src/data"
	"src/tcp"
)

func checkScreen() {
	//输出学生屏幕
	fmt.Println("当前在线学生屏幕信息如下")
	fmt.Println("学号\t姓名\t屏幕未改变时间")
	for userId := range userMgr {
		fmt.Printf("%s\t%s\t%d分钟\n", userMgr[userId].userData.UserId,
			userMgr[userId].userData.UserName, userMgr[userId].screenUnchangeTime)
	}
	//输出完成
	fmt.Println("您可以输入“学号”查看学生当前屏幕内容")
	fmt.Println("输入“.video”进入视频模式")
	var id string
	fmt.Scanf("%s\n", &id)
	if id != "" {
		if id == ".video" {
			fmt.Println("进入视频模式 请输入“学号”")
			fmt.Scanf("%s\n", &id)
			if id != "" {
				_, flag := userMgr[id]
				if flag {
					var msg data.Message
					msg.Type = data.ScreenVideoStartType
					dataByte, _ := json.Marshal(msg)
					tcp.WritePkg(userMgr[id].conn, dataByte)
					fmt.Println("已发送查看请求 输入任意内容可退出")
					exec.Command("cmd", "/c", "ScreenShot\\viewScreenVideo.html").Run()
					fmt.Scanf("%s\n", &id)
					msg.Type = data.ScreenVideoStopType
					dataByte, _ = json.Marshal(msg)
					tcp.WritePkg(userMgr[id].conn, dataByte)
				} else {
					fmt.Println("学号输入错误或学生未上线")
				}
			}
		} else {
			_, flag := userMgr[id]
			if flag {
				var msg data.Message
				msg.Type = data.ScreenShotGetType
				dataByte, _ := json.Marshal(msg)
				tcp.WritePkg(userMgr[id].conn, dataByte)
				fmt.Println("已发送查看请求")
			} else {
				fmt.Println("学号输入错误或学生未上线")
			}
		}
	}
}

func viewScreenShot(user data.User, screenShotRes data.ScreenShotRes) {
	dataByte, _ := base64.StdEncoding.DecodeString(screenShotRes.Img)
	os.Mkdir("ScreenShot", 0777)
	ioutil.WriteFile("ScreenShot/ScreenShot_"+user.UserId+".png", dataByte, 0777)
	exec.Command("cmd", "/c", "ScreenShot\\ScreenShot_"+user.UserId+".png").Run()
	fmt.Printf("收到屏幕截图：ScreenShot_" + user.UserId + ".png\n")
}

func viewScreenVideo(screenVideoRes data.ScreenVideoRes) {
	dataByte, _ := base64.StdEncoding.DecodeString(screenVideoRes.Img)
	os.Mkdir("ScreenShot", 0777)
	ioutil.WriteFile("ScreenShot/ScreenVideo.png", dataByte, 0777)
}
