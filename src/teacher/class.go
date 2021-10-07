package main

import (
	"encoding/json"
	"fmt"
	"net"
	"src/config"
	"src/data"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/garyburd/redigo/redis"
)

var class data.Class
var classData data.ClassData

//添加新课程
func addClass() {
	//读入课程名
	color.Cyan("请输入要添加的课程名：")
	fmt.Scanf("%s\n", &class.ClassName)

	rconn := RconnPool.Get()
	defer rconn.Close()
	rconn.Do("RPUSH", "class", class.ClassName)

	//切换到该课程的数据库
	r, _ := redis.Int(rconn.Do("LLEN", "class"))
	rconn.Do("select", r)

	//写入课程初始数据
	class.ClassId = r
	class.ClassNo = 0
	class.WorkNo = 0
	dataByte, _ := json.Marshal(class)
	rconn.Do("set", "classData", string(dataByte))

	//导入学生
	color.Cyan("请输入参加这门课程的学生数量:")
	var num int
	fmt.Scanf("%d\n", &num)
	color.Cyan("请按照“[学号] [密码] [姓名]”的格式导入学生数据\n")
	color.Cyan("如：S123 P456789 张三\n")
	for i := 0; i < num; i++ {
		var newUser data.User
		fmt.Scanf("%s %s %s\n", &newUser.UserId, &newUser.UserPwd, &newUser.UserName)
		userAdd(newUser)
	}
	color.Cyan("导入完成 课程已添加\n")

	//切换回主数据库
	rconn.Do("select", 0)
	class.ClassId = 0
}

//选择课程
func selectClass() {
	//连接数据库
	rconn := RconnPool.Get()
	defer rconn.Close()

	//获取课程总数
	r, _ := redis.Int(rconn.Do("LLEN", "class"))

	//输出课程
	color.Cyan("课程ID \t 课程名")
	for i := 0; i < r; i++ {
		class.ClassName, _ = redis.String(rconn.Do("LINDEX", "class", i))
		color.Cyan("%d \t %s \n", i+1, class.ClassName)
	}
	color.Cyan("请选择（1-%d）：\n", r)
	var key int
	fmt.Scanf("%d\n", &key)

	//切换课程数据库
	if key > 0 && key <= r {
		class.ClassId = key
		rconn.Do("select", class.ClassId)
		//读取课程数据
		dataByte, _ := redis.Bytes(rconn.Do("get", "classData"))
		json.Unmarshal(dataByte, &class)
		color.Cyan("您选择的课程是：%s\n", class.ClassName)
		color.Cyan("该课程已上%d节 有%d次作业\n", class.ClassNo, class.WorkNo)
		ShowClassMenu()
	}
}

//进程课堂
func enterClass() {
	//开始监听端口
	listen, err := net.Listen("tcp", "0.0.0.0:"+config.SerPoint)
	if err != nil {
		errorReport(err)
		return
	}
	defer listen.Close()

	//课程次数+1
	class.ClassNo++

	rconn := RconnPool.Get()
	defer rconn.Close()

	dataByte, _ := json.Marshal(class)
	rconn.Do("set", "classData", string(dataByte))

	classData.ClassNo = class.ClassNo
	dataByte, _ = json.Marshal(classData)
	rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), "classData", string(dataByte))

	go EnterClassMenu()
	//等待客户端来链接服务器
	for {
		conn, err := listen.Accept()
		if err != nil {
			errorReport(err)
		}
		//一旦链接成功，则启动一个线程和客户端保持通讯
		go process(conn)
	}
}

//开始上课
func beginClass() {
	sendMsg("同学们开始上课啦~")
	//记录开始时间 用来判断学生是否迟到
	classData.BeginTime = time.Now()
	classData.EndTime = time.Now().Add(time.Hour)
	dataByte, _ := json.Marshal(classData)

	rconn := RconnPool.Get()
	defer rconn.Close()
	rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), "classData", string(dataByte))

	BeginClassMenu()

}

//下课
func endClass() {
	sendMsg("下课啦 请同学们尽快退出登录~")
	//记录下课时间 用来判断学生是否早退等
	classData.EndTime = time.Now()
	dataByte, _ := json.Marshal(classData)

	rconn := RconnPool.Get()
	defer rconn.Close()
	rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), "classData", string(dataByte))

	EndClassMenu()
}

//查询课程数据
func checkClassData() {
	if class.ClassNo == 0 {
		color.Cyan("当前没有课程数据\n")
		return
	}

	color.Cyan("当前课程有%d个数据 您想看哪节课?\n", class.ClassNo)
	color.Cyan("请输入（1-%d）：\n", class.ClassNo)
	var key int
	fmt.Scanf("%d\n", &key)
	if key < 1 || key > class.ClassNo {
		color.Cyan("没有这节课的数据\n")
		return
	}

	rconn := RconnPool.Get()
	defer rconn.Close()
	res, _ := redis.String(rconn.Do("hget", "class"+strconv.Itoa(key), "classData"))
	var cData data.ClassData
	var uData data.User
	var ucData data.UserClassData
	json.Unmarshal([]byte(res), &cData)
	color.Cyan("您选择了第%d节课 持续时长%d分钟\n", cData.ClassNo, int(cData.EndTime.Sub(cData.BeginTime).Minutes()))
	color.Cyan("开始时间：%s 结束时间：%s\n",
		cData.BeginTime.Format("2006-01-02 15:04:05"), cData.EndTime.Format("2006-01-02 15:04:05"))
	color.Cyan("状态说明：1正常 2早退 3迟到 4迟到早退 5缺勤")
	color.Cyan("学号\t姓名\t状态\t加入时间\t\t离开时间\t\t违规次数\t机位")
	uDatas, _ := redis.Strings(rconn.Do("hkeys", "userData"))
	for uId := range uDatas {
		res, _ = redis.String(rconn.Do("hget", "userData", uDatas[uId]))
		json.Unmarshal([]byte(res), &uData)
		res, err := redis.String(rconn.Do("hget", "class"+strconv.Itoa(key), uDatas[uId]))
		if err == redis.ErrNil {
			//课程中没有这个学生的数据 缺勤!
			fmt.Printf("%s\t%s\t5\n", uData.UserId, uData.UserName)
		} else {
			json.Unmarshal([]byte(res), &ucData)
			if ucData.ClassStatus == 0 {
				//等于0 状态未更新
				if ucData.JoinTime.Before(cData.BeginTime) {
					//没有迟到
					if ucData.LeaveTime.After(cData.EndTime) {
						//没有早退
						ucData.ClassStatus = data.ClassStatus_Normal
					} else {
						//早退
						ucData.ClassStatus = data.ClassStatus_LeaveEarly
					}
				} else {
					//迟到了
					if ucData.LeaveTime.After(cData.EndTime) {
						//没早退
						ucData.ClassStatus = data.ClassStatus_JoinLate
					} else {
						ucData.ClassStatus = data.ClassStatus_LEJL
					}
				}
				//把状态记录到数据库
				dataByte, _ := json.Marshal(ucData)
				rconn.Do("hset", "class"+strconv.Itoa(key), uDatas[uId], string(dataByte))
			}
			fmt.Printf("%s\t%s\t%d\t%s\t%s\t%d\t%s\n", uData.UserId, uData.UserName, ucData.ClassStatus,
				ucData.JoinTime.Format("2006-01-02 15:04:05"), ucData.LeaveTime.Format("2006-01-02 15:04:05"),
				ucData.Violate, ucData.Seat)
		}
	}
	//输出学生数据完成
	color.Cyan("您可以输入“[学号] [状态]”进行修改 输入“0 0”返回菜单")
	var id string
	var status int
	for {
		fmt.Scanf("%s %d\n", &id, &status)
		if id == "0" && status == 0 {
			return
		} else if status < 1 || status > 5 {
			color.Cyan("输入有误，请重新输入")
		} else {
			res, err := redis.String(rconn.Do("hget", "class"+strconv.Itoa(key), id))
			if err == redis.ErrNil {
				//没数据 直接提交一个新的
				ucData.UserId = id
				ucData.ClassStatus = status
				ucData.JoinTime = time.Now()
				ucData.LeaveTime = time.Now()
				ucData.Violate = 0
				ucData.Seat = "add by teacher"
			} else {
				//有数据 修改
				json.Unmarshal([]byte(res), &ucData)
				ucData.ClassStatus = status
			}
			//记录到数据库
			dataByte, _ := json.Marshal(ucData)
			rconn.Do("hset", "class"+strconv.Itoa(key), id, string(dataByte))
		}
	}
}
