package main

import (
	"encoding/json"
	"fmt"
	"net"
	"src/data"
	"src/tcp"
	"strconv"
	"time"
)

//处理和客户端的通讯
func process(conn net.Conn) {
	defer conn.Close()
	var user data.User
	var userClassData data.UserClassData
	var userProcess UserProcess
	for {
		time.Sleep(time.Millisecond * 100)
		//延迟一下 确保ReadPkg能读到完整数据
		msg, err := tcp.ReadPkg(conn)
		if err != nil {
			//出现错误 掉线了
			if user.UserName != "" {
				fmt.Printf("[L][%s]学生%s异常掉线\n", user.UserId, user.UserName)
				userClassData.LeaveTime = time.Now()
				dataByte, _ := json.Marshal(userClassData)
				Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), user.UserId, string(dataByte))
				DelOnlineUser(user.UserId)
			}
			return
		}
		if msg.Type == "" {
			return
		}
		switch msg.Type {
		case data.LoginMesType: //处理登录
			var loginMes data.LoginMes
			json.Unmarshal([]byte(msg.Data), &loginMes)
			var loginResMes data.LoginResMes
			loginResMes.Result, loginResMes.Username = userLogin(loginMes)
			loginResMes.ClassName = class.ClassName
			if loginResMes.Result == data.Logout_Success {
				user.UserId = loginMes.UserId
				user.UserName = loginResMes.Username
				userClassData.Seat = loginMes.Seat
				userClassData.UserId = loginMes.UserId
				userClassData.JoinTime = time.Now()
				userClassData.ClassStatus = 0
				dataByte, _ := json.Marshal(userClassData)
				Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), user.UserId, string(dataByte))
				userProcess.userData = user
				userProcess.conn = conn
				AddOnlineUser(&userProcess)
				onlineNum := GetOnlineUserNum()
				fmt.Printf("[L][%s]学生%s已登录 机位为%s 当前在线：%d\n", loginMes.UserId, user.UserName, loginMes.Seat, onlineNum)
			}

			//使用json序列化
			dataByte, _ := json.Marshal(loginResMes)
			msg.Type = data.LoginResMesType
			msg.Data = string(dataByte)
			dataByte, _ = json.Marshal(msg)

			//发送消息给客户
			tcp.WritePkg(conn, []byte(dataByte))

			if loginResMes.Result == data.Logout_Success {
				go sendBlockList(conn)
			}
		case data.LogoutMesType: //处理注销
			var logoutMes data.LogoutMes
			json.Unmarshal([]byte(msg.Data), &logoutMes)
			var logoutResMes data.LogoutResMes
			if logoutMes.UserId == user.UserId {
				logoutResMes.Result = data.Logout_Success
				userClassData.LeaveTime = time.Now()
				dataByte, _ := json.Marshal(userClassData)
				Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), user.UserId, string(dataByte))
				DelOnlineUser(user.UserId)
				fmt.Printf("[L][%s]学生%s已离开\n", user.UserId, user.UserName)
				//使用json序列化
				dataByte, _ = json.Marshal(logoutResMes)
				msg.Type = data.LogoutResMesType
				msg.Data = string(dataByte)
				dataByte, _ = json.Marshal(msg)

				//发送消息给客户
				tcp.WritePkg(conn, []byte(dataByte))
				return
			} else {
				fmt.Printf("[W]注销消息出现错误 %s %s\n", logoutMes.UserId, user.UserId)
			}
		case data.ChatMesType: //处理聊天
			var chatMes data.ChatMes
			json.Unmarshal([]byte(msg.Data), &chatMes)
			sendResMsg(user, chatMes)
		case data.ChatPMesType: //处理私聊
			var chatPMes data.ChatPMes
			json.Unmarshal([]byte(msg.Data), &chatPMes)
			sendPResMsg(user, chatPMes)
		case data.QueMesType:
			var queMes data.QueMes
			json.Unmarshal([]byte(msg.Data), &queMes)
			askQueRes(user, queMes, userClassData.Seat)
		case data.WorkAllMesType:
			sendWorkAll(conn)
		case data.WorkMesType:
			var workMes data.WorkMes
			json.Unmarshal([]byte(msg.Data), &workMes)
			sendWorkData(user, conn, workMes)
		case data.WorkSubMesType:
			var workSubMes data.WorkSubMes
			json.Unmarshal([]byte(msg.Data), &workSubMes)
			sendWorkSub(user, conn, workSubMes)
		case data.ScreenReportType:
			var screenReport data.ScreenReport
			json.Unmarshal([]byte(msg.Data), &screenReport)
			userMgr[user.UserId].screenUnchangeTime = screenReport.UnchangeTime
			if screenReport.UnchangeTime > 2 {
				fmt.Printf("[I][%s]学生%s屏幕内容已连续%d分钟未改变\n", user.UserId, user.UserName, screenReport.UnchangeTime)
			}
		case data.ScreenShotResType:
			var screenShotRes data.ScreenShotRes
			json.Unmarshal([]byte(msg.Data), &screenShotRes)
			viewScreenShot(user, screenShotRes)
		case data.ScreenVideoResType:
			var screenVideoRes data.ScreenVideoRes
			json.Unmarshal([]byte(msg.Data), &screenVideoRes)
			viewScreenVideo(screenVideoRes)
			conn.Close()
			return
		default:
			fmt.Printf("[W]连接%s的消息类型为%s 无法处理\n", conn.RemoteAddr().String(), msg.Type)
			return
		}
	}
}
