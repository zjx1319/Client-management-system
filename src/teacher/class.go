package main

import (
	"encoding/json"
	"fmt"
	"net"
	"src/config"
	"src/data"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

func addClass() {
	//读入课程名
	fmt.Print("请输入课程名：")
	fmt.Scanf("%s\n", &class.ClassName)
	fmt.Printf("课程名为：%s\n", class.ClassName)
	Rconn.Do("RPUSH", "class", class.ClassName)

	//切换到该课程的数据库
	r, _ := redis.Int(Rconn.Do("LLEN", "class"))
	Rconn.Do("select", r)

	//写入课程初始数据
	class.ClassId = r
	class.ClassNo = 0
	class.WorkNo = 0
	dataByte, _ := json.Marshal(class)
	Rconn.Do("set", "classData", string(dataByte))

	//导入学生
	fmt.Println("导入学生")
	fmt.Print("请输入要导入的学生数量:")
	var num int
	fmt.Scanf("%d\n", &num)
	fmt.Print("请按照“学号 密码 姓名”的格式导入学生数据\n")
	for i := 0; i < num; i++ {
		var newUser data.User
		fmt.Scanf("%s %s %s\n", &newUser.UserId, &newUser.UserPwd, &newUser.UserName)
		userAdd(newUser)
	}
	fmt.Println("导入完成")

	//切换回主数据库
	Rconn.Do("select", 0)
}

func selectClass() {

	//获取课程总数
	r, err := redis.Int(Rconn.Do("LLEN", "class"))
	if err != nil {
		fmt.Println("LLEN err=", err)
		return
	}

	//输出课程
	fmt.Println("课程ID \t 课程名")
	for i := 0; i < r; i++ {
		class.ClassName, err = redis.String(Rconn.Do("LINDEX", "class", i))
		if err != nil {
			fmt.Println("LINDEX err=", err)
			return
		}
		fmt.Printf("%d \t %s \n", i+1, class.ClassName)
	}
	fmt.Println("输入其他数字返回主菜单")
	fmt.Print("请选择:")
	fmt.Scanf("%d\n", &class.ClassId)

	//切换课程数据库
	if class.ClassId > 0 && class.ClassId <= r {
		_, err = redis.String(Rconn.Do("select", class.ClassId))
		if err != nil {
			fmt.Println("select err=", err)
			return
		}
		//读取课程数据
		b, _ := redis.Bytes(Rconn.Do("get", "classData"))
		json.Unmarshal(b, &class)
		fmt.Printf("您选择的课程是：%s\n", class.ClassName)
		fmt.Printf("该课程已上%d节 有%d次作业\n", class.ClassNo, class.WorkNo)
		ShowClassMenu()
	}
}

func enterClass() {
	//课程次数+1 数据库信息处理
	class.ClassNo++
	dataByte, _ := json.Marshal(class)
	Rconn.Do("set", "classData", string(dataByte))

	classData.ClassNo = class.ClassNo

	dataByte, _ = json.Marshal(classData)
	Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), "classData", string(dataByte))

	//提示信息
	listen, err := net.Listen("tcp", "0.0.0.0:"+config.SerPoint)
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()

	go EnterClassMenu()
	//等待客户端来链接服务器
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//一旦链接成功，则启动一个线程和客户端保持通讯
		go process(conn)
	}
}

func beginClass() {
	//记录开始时间 用来判断学生是否迟到
	classData.BeginTime = time.Now()
	dataByte, _ := json.Marshal(classData)
	Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), "classData", string(dataByte))
	BeginClassMenu()
}

func endClass() {
	//记录下课时间 用来判断学生是否早退等
	classData.EndTime = time.Now()
	dataByte, _ := json.Marshal(classData)
	Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), "classData", string(dataByte))
	EndClassMenu()
}

func checkClassData() {
	var key int
	if class.ClassNo == 0 {
		fmt.Println("当前没有课程数据")
		return
	}
	fmt.Printf("当前课程有%d个数据 您想看哪节课?\n", class.ClassNo)
	fmt.Printf("请输入（1-%d）：\n", class.ClassNo)
	fmt.Scanf("%d\n", &key)
	if key < 1 || key > class.ClassNo {
		fmt.Println("没有这节课的数据")
		return
	}

	res, _ := redis.String(Rconn.Do("hget", "class"+strconv.Itoa(key), "classData"))
	var cData data.ClassData
	var uData data.User
	var ucData data.UserClassData
	json.Unmarshal([]byte(res), &cData)
	fmt.Printf("您选择了第%d节课 持续时长%d分钟\n", cData.ClassNo, int(cData.EndTime.Sub(cData.BeginTime).Minutes()))
	fmt.Printf("开始时间：%s 结束时间：%s\n",
		cData.BeginTime.Format("2006-01-02 15:04:05"), cData.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Println("状态说明：1正常 2早退 3迟到 4迟到早退 5缺勤")
	fmt.Println("学号\t姓名\t状态\t加入时间\t\t离开时间\t\t违规次数\t机位")
	uDatas, _ := redis.Strings(Rconn.Do("hkeys", "userData"))
	for uId := range uDatas {
		res, _ = redis.String(Rconn.Do("hget", "userData", uDatas[uId]))
		json.Unmarshal([]byte(res), &uData)
		res, err := redis.String(Rconn.Do("hget", "class"+strconv.Itoa(key), uDatas[uId]))
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
				Rconn.Do("hset", "class"+strconv.Itoa(key), uDatas[uId], string(dataByte))
			}
			fmt.Printf("%s\t%s\t%d\t%s\t%s\t%d\t%s\n", uData.UserId, uData.UserName, ucData.ClassStatus,
				ucData.JoinTime.Format("2006-01-02 15:04:05"), ucData.LeaveTime.Format("2006-01-02 15:04:05"),
				ucData.Violate, ucData.Seat)
		}
	}
	//输出学生数据完成
	fmt.Println("您可以输入“学号 状态”进行修改 输入“0 0”返回菜单")
	var id string
	var status int
	for {
		fmt.Scanf("%s %d\n", &id, &status)
		if id == "0" && status == 0 {
			return
		} else if status < 1 || status > 5 {
			fmt.Println("您输入的状态有误 请重新输入")
		} else {
			res, err := redis.String(Rconn.Do("hget", "class"+strconv.Itoa(key), id))
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
			Rconn.Do("hset", "class"+strconv.Itoa(key), id, string(dataByte))
		}
	}
}
