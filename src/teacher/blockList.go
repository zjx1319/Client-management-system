package main

import (
	"encoding/json"
	"net"
	"src/data"
	"src/tcp"
	"strconv"

	"github.com/fatih/color"
)

var BlockListProcess []string
var BlockListWeb []string

func init() {
	BlockListProcess = []string{"steam.exe", "wegame.exe"}
	BlockListWeb = []string{"4399.com", "www.4399.com"}
}

func sendBlockList(conn net.Conn) {
	var msg data.Message
	msg.Type = data.BlockListProcessType
	var blockListProcess data.BlockListProcess
	blockListProcess.List = BlockListProcess
	dataByte, _ := json.Marshal(blockListProcess)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)
	tcp.WritePkg(conn, dataByte)

	msg.Type = data.BlockListWebType
	var blockListWeb data.BlockListWeb
	blockListWeb.List = BlockListWeb
	dataByte, _ = json.Marshal(blockListWeb)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)
	tcp.WritePkg(conn, dataByte)
}

func BlockListDeal(ucData *data.UserClassData, user *data.User, behavior string) {
	color.Yellow("[Violation][%s]学生%s违规操作：%s\n", user.UserId, user.UserName, behavior)
	ucData.Violate++
	dataByte, _ := json.Marshal(ucData)
	rconn := RconnPool.Get()
	defer rconn.Close()
	rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), user.UserId, string(dataByte))
}
