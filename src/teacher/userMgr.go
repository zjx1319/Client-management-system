package main

import (
	"net"
	"src/data"
)

var userMgr map[string]*UserProcess

type UserProcess struct {
	userData           data.User
	screenUnchangeTime int
	conn               net.Conn
}

func init() {
	userMgr = make(map[string]*UserProcess)
}

//添加
func AddOnlineUser(user *UserProcess) {
	userMgr[user.userData.UserId] = user
}

//删除
func DelOnlineUser(userId string) {
	delete(userMgr, userId)
}

//根据id返回对应的值
func GetOnlineUserById(userId string) (up *UserProcess) {
	up = userMgr[userId]
	return
}

//获取在线用户数量
func GetOnlineUserNum() (num int) {
	num = len(userMgr)
	return
}
