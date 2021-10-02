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
	for {
		msg, err := tcp.ReadPkg(conn)
		if err != nil {
			//出现错误 掉线了
			if user.UserName != "" {
				fmt.Printf("【%s】学生%s异常掉线\n", user.UserId, user.UserName)
				userClassData.LeaveTime = time.Now()
				dataByte, _ := json.Marshal(userClassData)
				Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), user.UserId, string(dataByte))
				classData.StudentNum--
			}
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
				classData.StudentNum++
				fmt.Printf("【%s】学生%s已登录 机位为%s 当前在线：%d\n", loginMes.UserId, user.UserName, loginMes.Seat, classData.StudentNum)
			}

			//使用json序列化
			dataByte, _ := json.Marshal(loginResMes)
			msg.Type = data.LoginResMesType
			msg.Data = string(dataByte)
			dataByte, _ = json.Marshal(msg)

			//发送消息给客户
			tcp.WritePkg(conn, []byte(dataByte))
		case data.LogoutMesType: //处理注销
			var logoutMes data.LogoutMes
			json.Unmarshal([]byte(msg.Data), &logoutMes)
			var logoutResMes data.LogoutResMes
			if logoutMes.UserId == user.UserId {
				logoutResMes.Result = data.Logout_Success
				userClassData.LeaveTime = time.Now()
				dataByte, _ := json.Marshal(userClassData)
				Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), user.UserId, string(dataByte))
				classData.StudentNum--
				fmt.Printf("【%s】学生%s已离开\n", user.UserId, user.UserName)
				//使用json序列化
				dataByte, _ = json.Marshal(logoutResMes)
				msg.Type = data.LogoutResMesType
				msg.Data = string(dataByte)
				dataByte, _ = json.Marshal(msg)

				//发送消息给客户
				tcp.WritePkg(conn, []byte(dataByte))
				return
			} else {
				fmt.Printf("注销消息出现错误 %s %s\n", logoutMes.UserId, user.UserId)
			}
		default:
			fmt.Printf("消息类型为%s 无法处理\n", msg.Type)
			return
		}
	}
}
