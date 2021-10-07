package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"src/data"
	"src/tcp"
	"time"

	"github.com/fatih/color"
)

var flag bool

func getWorkAll() {
	var msg data.Message
	var dataByte []byte
	flag = true
	msg.Type = data.WorkAllMesType
	dataByte, _ = json.Marshal(msg)
	//发送数据
	tcp.WritePkg(conn, dataByte)
	for flag {
		time.Sleep(time.Second)
		//等待教师端返回数据
	}
}

func getWorkData(num int) {
	var workMes data.WorkMes
	var msg data.Message
	var dataByte []byte
	if num == 0 {
		color.Cyan("当前没有作业数据\n")
		flag = false
		return
	}
	color.Cyan("请输入你想查看的作业Id(1-%d):\n", num)
	fmt.Scanf("%d\n", &workMes.Id)
	if workMes.Id < 1 || workMes.Id > num {
		color.Cyan("没有这个作业数据\n")
		flag = false
		return
	}
	msg.Type = data.WorkMesType
	dataByte, _ = json.Marshal(workMes)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)
	//发送数据
	tcp.WritePkg(conn, dataByte)
}

func workDataRes(workResMes data.WorkResMes) {
	color.Cyan("你查询的作业Id为%d\n", workResMes.Id)
	switch workResMes.Type {
	case data.Work_Objective:
		color.Cyan("题目类型：客观题 ")
	case data.Work_Subjective:
		color.Cyan("题目类型：主观题 ")
	case data.Work_Files:
		color.Cyan("题目类型：文件题 ")
	default:
		color.Cyan("题目类型：未知 ")
	}
	color.Cyan("满分：%d 截止时间：%s\n", workResMes.FullScore, workResMes.DeadLine.Format("2006-01-02 15:04:05"))
	color.Cyan("题目：%s\n", workResMes.Question)
	if workResMes.Answer != "" {
		//已提交
		color.Cyan("你的答案：%s\n得分：%d 提交时间：%s\n", workResMes.Answer, workResMes.Score,
			workResMes.SubmitTime.Format("2006-01-02 15:04:05"))
	} else {
		if workResMes.Type != data.Work_Files {
			color.Cyan("要提交答案请直接输入 不输入直接回车可返回主菜单\n")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			if input.Text() != "" {
				var workSubMes data.WorkSubMes
				var msg data.Message
				var dataByte []byte
				msg.Type = data.WorkSubMesType
				workSubMes.Id = workResMes.Id
				workSubMes.Answer = input.Text()
				dataByte, _ = json.Marshal(workSubMes)
				msg.Data = string(dataByte)
				dataByte, _ = json.Marshal(msg)
				//发送数据
				tcp.WritePkg(conn, dataByte)
				color.Cyan("作业已提交\n")
			}
		} else {
			//文件提交
			color.Cyan("文件题要提交答案请输入文件路径（不带引号 不超过5M） 不输入直接回车可返回主菜单\n")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			if input.Text() != "" {
				var workSubMes data.WorkSubMes
				var msg data.Message
				var dataByte []byte
				msg.Type = data.WorkSubMesType
				workSubMes.Id = workResMes.Id
				//读取文件
				f, err := ioutil.ReadFile(input.Text())
				if err != nil {
					color.Cyan("文件读取出错 提交失败\n")
				} else if len(f) > 5242880 {
					color.Cyan("文件过大 请提交5M以下的文件\n")
				} else {
					var fileMes data.FileMes
					fileMes.FileName = filepath.Base(input.Text())
					fileMes.Data = base64.StdEncoding.EncodeToString(f)
					dataByte, _ = json.Marshal(fileMes)
					workSubMes.Answer = string(dataByte)
					dataByte, _ = json.Marshal(workSubMes)
					msg.Data = string(dataByte)
					dataByte, _ = json.Marshal(msg)
					//发送数据
					tcp.WritePkg(conn, dataByte)
					color.Cyan("作业已提交\n")
				}

			}
		}
	}
	flag = false
}

func subWorkRes(workSubResMes data.WorkSubResMes) {
	fmt.Printf("作业%d提交成功 得分%d\n", workSubResMes.Id, workSubResMes.Score)
	if workSubResMes.Score == 0 {
		fmt.Println("得分为0可能是提交超时或未评分或客观题做错啦")
	}
}
