package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"src/data"
	"src/tcp"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

func addWork() {
	var workData data.WorkData

	fmt.Println("请选择题目类型（1是客观题 2是主观题 3是文件题）：")
	fmt.Scanf("%d\n", &workData.Type)
	if workData.Type < 1 || workData.Type > 3 {
		fmt.Println("您输入的类型错误！")
		return
	}
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("请输入题目：")
	input.Scan()
	workData.Question = input.Text()
	fmt.Println("请输入参考答案（客观题会自动打分）：")
	input.Scan()
	workData.StdAnswer = input.Text()
	fmt.Println("请输入题目分值：")
	fmt.Scanf("%d\n", &workData.FullScore)
	fmt.Println("您想在多久后截止提交？（输入整数 单位分钟）")
	var ddlm int
	fmt.Scanf("%d\n", &ddlm)
	timeAdd, _ := time.ParseDuration(strconv.Itoa(ddlm) + "m")
	workData.DeadLine = time.Now().Add(timeAdd)

	//作业次数+1 数据库信息处理
	class.WorkNo++
	dataByte, _ := json.Marshal(class)
	Rconn.Do("set", "classData", string(dataByte))

	workData.Id = class.WorkNo

	dataByte, _ = json.Marshal(workData)
	Rconn.Do("hset", "work"+strconv.Itoa(workData.Id), "workData", string(dataByte))

	//作业添加成功
	sendMsg("新的作业已发布，Id为" + strconv.Itoa(workData.Id) + " 请在" +
		workData.DeadLine.Format("15:04:05") + "前提交")
}

func checkWorkData() {
	var key int
	if class.WorkNo == 0 {
		fmt.Println("当前没有作业数据")
		return
	}
	fmt.Printf("当前课程有%d个作业数据 您想看哪个作业?\n", class.WorkNo)
	fmt.Printf("请输入（1-%d）：\n", class.WorkNo)
	fmt.Scanf("%d\n", &key)
	if key < 1 || key > class.WorkNo {
		fmt.Println("没有这个作业的数据")
		return
	}

	res, _ := redis.String(Rconn.Do("hget", "work"+strconv.Itoa(key), "workData"))
	var wData data.WorkData
	var uData data.User
	var uwData data.UserWorkData
	json.Unmarshal([]byte(res), &wData)
	//输出题目信息
	fmt.Printf("你查询的作业Id为%d\n", wData.Id)
	switch wData.Type {
	case data.Work_Objective:
		fmt.Print("题目类型：客观题 ")
	case data.Work_Subjective:
		fmt.Print("题目类型：主观题 ")
	case data.Work_Files:
		fmt.Print("题目类型：文件题 ")
	default:
		fmt.Print("题目类型：未知 ") //Error
	}
	fmt.Printf("满分：%d 截止时间：%s\n", wData.FullScore, wData.DeadLine.Format("2006-01-02 15:04:05"))
	fmt.Printf("题目：%s\n", wData.Question)
	fmt.Printf("参考答案：%s\n", wData.StdAnswer)
	//输出学生答题信息
	fmt.Println("学号\t姓名\t得分\t提交时间\t\t答案")
	uDatas, _ := redis.Strings(Rconn.Do("hkeys", "userData"))
	for uId := range uDatas {
		res, _ = redis.String(Rconn.Do("hget", "userData", uDatas[uId]))
		json.Unmarshal([]byte(res), &uData)
		res, err := redis.String(Rconn.Do("hget", "work"+strconv.Itoa(key), uDatas[uId]))
		if err == redis.ErrNil {
			//课程中没有这个学生的数据 0分!
			fmt.Printf("%s\t%s\t0\t未提交\n", uData.UserId, uData.UserName)
		} else {
			json.Unmarshal([]byte(res), &uwData)
			fmt.Printf("%s\t%s\t%d\t%s\t%s\n", uData.UserId, uData.UserName, uwData.Score,
				uwData.SubmitTime.Format("2006-01-02 15:04:05"), uwData.Answer)
		}
	}
	//输出学生数据完成
	fmt.Println("您可以输入“学号 分数”进行批改 输入“0 0”返回菜单")
	var id string
	var score int
	for {
		fmt.Scanf("%s %d\n", &id, &score)
		if id == "0" && score == 0 {
			return
		} else {
			res, err := redis.String(Rconn.Do("hget", "work"+strconv.Itoa(key), id))
			if err == redis.ErrNil {
				//没数据 直接提交一个新的
				uwData.UserId = id
				uwData.Score = score
				uwData.SubmitTime = time.Now()
				uwData.Answer = "mod by teacher"
			} else {
				//有数据 修改
				json.Unmarshal([]byte(res), &uwData)
				uwData.Score = score
			}
			//记录到数据库
			dataByte, _ := json.Marshal(uwData)
			Rconn.Do("hset", "work"+strconv.Itoa(key), id, string(dataByte))
		}
	}
}

func sendWorkAll(conn net.Conn) {
	var msg data.Message
	var workAllResMes data.WorkAllResMes
	var dataByte []byte
	workAllResMes.Num = class.WorkNo
	msg.Type = data.WorkAllResMesType
	dataByte, _ = json.Marshal(workAllResMes)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)
	tcp.WritePkg(conn, dataByte)
}

func sendWorkData(user data.User, conn net.Conn, workMes data.WorkMes) {
	var msg data.Message
	var workResMes data.WorkResMes
	var dataByte []byte

	//获取课程数据
	workResMes.Id = workMes.Id
	res, _ := redis.String(Rconn.Do("hget", "work"+strconv.Itoa(workResMes.Id), "workData"))
	var wData data.WorkData
	json.Unmarshal([]byte(res), &wData)
	workResMes.Type = wData.Type
	workResMes.Question = wData.Question
	workResMes.DeadLine = wData.DeadLine
	workResMes.FullScore = wData.FullScore

	var uwData data.UserWorkData
	res, err := redis.String(Rconn.Do("hget", "work"+strconv.Itoa(workResMes.Id), user.UserId))
	if err == redis.ErrNil {
		//没数据
		workResMes.Answer = ""
	} else {
		//有数据
		json.Unmarshal([]byte(res), &uwData)
		workResMes.Answer = uwData.Answer
		workResMes.SubmitTime = uwData.SubmitTime
		workResMes.Score = uwData.Score
	}

	//发送数据
	msg.Type = data.WorkResMesType
	dataByte, _ = json.Marshal(workResMes)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)
	tcp.WritePkg(conn, dataByte)
}

func sendWorkSub(user data.User, conn net.Conn, workSubMes data.WorkSubMes) {
	var msg data.Message
	var workSubResMes data.WorkSubResMes
	var dataByte []byte
	var score int

	var uwData data.UserWorkData
	uwData.UserId = user.UserId
	uwData.Answer = workSubMes.Answer
	uwData.SubmitTime = time.Now()

	//获取课程数据
	workSubResMes.Id = workSubMes.Id
	res, _ := redis.String(Rconn.Do("hget", "work"+strconv.Itoa(workSubResMes.Id), "workData"))
	var wData data.WorkData
	json.Unmarshal([]byte(res), &wData)
	//文件题先写出文件并把答案改成文件目录
	if wData.Type == data.Work_Files {
		var fileMes data.FileMes
		json.Unmarshal([]byte(workSubMes.Answer), &fileMes)
		os.MkdirAll("work/"+strconv.Itoa(workSubResMes.Id)+"/"+user.UserId, 0777)
		fileName := "work/" + strconv.Itoa(workSubResMes.Id) + "/" + user.UserId + "/" + fileMes.FileName
		dataByte, _ = base64.StdEncoding.DecodeString(fileMes.Data)
		ioutil.WriteFile(fileName, dataByte, 0777)
		uwData.Answer = "work/" + strconv.Itoa(workSubResMes.Id) + "/" + user.UserId + "/" + fileMes.FileName
	}

	//如果是客观题自动改分
	if wData.Type == data.Work_Objective && wData.StdAnswer == workSubMes.Answer {
		score = wData.FullScore
	} else {
		score = 0
	}
	//判断提交是否超时
	if uwData.SubmitTime.After(wData.DeadLine) {
		score = 0
	}

	uwData.Score = score
	workSubResMes.Score = score

	//记录到数据库
	dataByte, _ = json.Marshal(uwData)
	Rconn.Do("hset", "work"+strconv.Itoa(workSubResMes.Id), uwData.UserId, string(dataByte))

	//发送数据
	msg.Type = data.WorkSubResMesType
	dataByte, _ = json.Marshal(workSubResMes)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)
	tcp.WritePkg(conn, dataByte)

	fmt.Printf("[W][%s]学生%s提交了作业%d：%s\n", uwData.UserId, user.UserName, workSubMes.Id, uwData.Answer)
}
