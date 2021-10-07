package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"src/data"
	"src/tcp"
	"strings"
	"time"

	"github.com/fatih/color"
)

var WebBlockList []string

func blockListCheckProcess(blockList []string) {
	var msg data.Message
	msg.Type = data.BlockListReportType
	for {
		for i := range blockList {
			if isProcessExist(blockList[i]) {
				var BlockListReport data.BlockListReport
				BlockListReport.Behavior = "黑名单进程：" + blockList[i]
				dataByte, _ := json.Marshal(BlockListReport)
				msg.Data = string(dataByte)
				dataByte, _ = json.Marshal(msg)
				//发送数据
				tcp.WritePkg(conn, dataByte)
				color.Yellow("发现黑名单进程：%s\n", blockList[i])
				exec.Command("taskkill", "/IM", blockList[i]).Run()
			}
		}
		time.Sleep(time.Second)
	}
}

//进程是否存在
func isProcessExist(processName string) bool {
	cmd := exec.Command("cmd", "/C", "tasklist")
	output, _ := cmd.Output()
	return strings.Contains(string(output), processName)
}

func blockListWeb(blockList []string) {
	WebBlockList = blockList
	for i := range WebBlockList {
		exec.Command("cmd", "/C", "echo 127.0.0.1      "+blockList[i]+" >> C:\\Windows\\System32\\drivers\\etc\\hosts").Run()
	}
	http.HandleFunc("/", blockListWebCheck)
	http.ListenAndServe("127.0.0.1:80", nil)
}

func blockListWebCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("上课期间请不要访问本网址!"))
	for i := range WebBlockList {
		if strings.Contains(r.RequestURI, "favicon.ico") && strings.Contains(r.Host, WebBlockList[i]) {
			var blockListReport data.BlockListReport
			var msg data.Message
			blockListReport.Behavior = "黑名单网站：" + WebBlockList[i]
			dataByte, _ := json.Marshal(blockListReport)
			msg.Data = string(dataByte)
			msg.Type = data.BlockListReportType
			dataByte, _ = json.Marshal(msg)
			//发送数据
			tcp.WritePkg(conn, dataByte)
			color.Yellow("发现黑名单网站：%s\n" + WebBlockList[i])
		}
	}
}
