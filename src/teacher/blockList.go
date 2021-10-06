package main

import (
	"encoding/json"
	"net"
	"src/data"
	"src/tcp"
)

var blockListProcess []string

func init() {
	blockListProcess = []string{"steam.exe", "wegame.exe"}
}
func sendBlockList(conn net.Conn) {
	var msg data.Message
	msg.Type = data.blockListProcessType
	var bolckListProcess data.BlockListProcess
	bolckListProcess.List = blockListProcess
	dataByte, _ := json.Marshal(bolckListProcess)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)
	tcp.WritePkg(conn, dataByte)
}
