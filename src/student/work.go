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
		fmt.Println("当前没有作业数据")
		flag = false
		return
	}
	fmt.Printf("请输入你想查看的作业Id(1-%d):\n", num)
	fmt.Scanf("%d\n", &workMes.Id)
	if workMes.Id < 1 || workMes.Id > num {
		fmt.Println("没有这个作业数据")
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
	fmt.Printf("你查询的作业Id为%d\n", workResMes.Id)
	switch workResMes.Type {
	case data.Work_Objective:
		fmt.Print("题目类型：客观题 ")
	case data.Work_Subjective:
		fmt.Print("题目类型：主观题 ")
	case data.Work_Files:
		fmt.Print("题目类型：文件题 ")
	default:
		fmt.Print("题目类型：未知 ") //Error
	}
	fmt.Printf("满分：%d 截止时间：%s\n", workResMes.FullScore, workResMes.DeadLine.Format("2006-01-02 15:04:05"))
	fmt.Printf("题目：%s\n", workResMes.Question)
	if workResMes.Answer != "" {
		//已提交
		fmt.Printf("你的答案：%s\n得分：%d 提交时间：%s\n", workResMes.Answer, workResMes.Score,
			workResMes.SubmitTime.Format("2006-01-02 15:04:05"))
	} else {
		if workResMes.Type != data.Work_Files {
			fmt.Println("要提交答案请直接输入 不输入直接回车可返回主菜单")
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
				fmt.Println("作业已提交")
			}
		} else {
			//文件提交
			fmt.Println("文件题要提交答案请输入文件路径（不带引号 不超过5M） 不输入直接回车可返回主菜单")
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
					fmt.Println("文件读取出错 提交失败")
				} else if len(f) > 5242880 {
					fmt.Println("文件过大 请提交5M以下的文件")
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
					fmt.Println("作业已提交")
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
