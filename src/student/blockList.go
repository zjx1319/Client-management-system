package main

import (
	"os/exec"
	"strings"
	"time"
)

func blockListCheckProcess(blockList []string) {
	for {
		for i := range blockList {
			if isProcessExist(blockList[i]) {
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
