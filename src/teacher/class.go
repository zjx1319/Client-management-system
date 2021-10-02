package main

import (
	"encoding/json"
	"fmt"
	"net"
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

func selectClass() (err error) {

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
	return
}

func enterClass() {
	//课程次数+1 数据库信息处理
	class.ClassNo++
	dataByte, _ := json.Marshal(class)
	Rconn.Do("set", "classData", string(dataByte))

	classData.ClassNo = class.ClassNo
	classData.StudentNum = 0

	dataByte, _ = json.Marshal(classData)
	Rconn.Do("hset", "class"+strconv.Itoa(class.ClassNo), "classData", string(dataByte))

	//提示信息
	fmt.Println("开始在8989端口监听....")
	listen, err := net.Listen("tcp", "0.0.0.0:8989")
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
