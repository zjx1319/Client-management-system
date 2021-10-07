package main

import (
	"encoding/json"
	"fmt"
	"src/data"
	"strconv"

	"github.com/fatih/color"
	"github.com/garyburd/redigo/redis"
)

func userAdd(newUser data.User) {
	rconn := RconnPool.Get()
	defer rconn.Close()
	data, _ := json.Marshal(newUser)
	rconn.Do("HSET", "userData", newUser.UserId, string(data))
}

func userLogin(loginMes data.LoginMes) (result string, userName string) {
	rconn := RconnPool.Get()
	defer rconn.Close()
	res, err := redis.String(rconn.Do("hget", "userData", loginMes.UserId))
	if err == redis.ErrNil { //表示在 users 哈希中，没有找到对应id
		result = data.Login_IDNotFound
		return
	}
	var user data.User
	json.Unmarshal([]byte(res), &user)
	if user.UserPwd != loginMes.UserPwd {
		result = data.Login_PwdError
		return
	}
	result = data.Login_Success
	userName = user.UserName
	return
}

func checkUserData() {
	color.Cyan("学号\t姓名\t正常到课\t作业次数\t平均分\t违规次数\n")

	rconn := RconnPool.Get()
	defer rconn.Close()

	var uData data.User
	var ucData data.UserClassData
	var uwData data.UserWorkData
	uDatas, _ := redis.Strings(rconn.Do("hkeys", "userData"))
	for uId := range uDatas {
		var classTimes, workTimes, score, violationTimes int
		res, _ := redis.String(rconn.Do("hget", "userData", uDatas[uId]))
		json.Unmarshal([]byte(res), &uData)
		for i := 1; i <= classData.ClassNo; i++ {
			res, err := redis.String(rconn.Do("hget", "class"+strconv.Itoa(i), uDatas[uId]))
			if err != redis.ErrNil {
				json.Unmarshal([]byte(res), &ucData)
				violationTimes = violationTimes + ucData.Violate
				if ucData.ClassStatus == 1 {
					classTimes++
				}
			}
		}
		for i := 1; i <= class.WorkNo; i++ {
			res, err := redis.String(rconn.Do("hget", "work"+strconv.Itoa(i), uDatas[uId]))
			if err != redis.ErrNil {
				json.Unmarshal([]byte(res), &uwData)
				workTimes++
				score = score + uwData.Score
			}
		}
		if workTimes == 0 {
			fmt.Printf("%s\t%s\t%d\t\t%d\t\t0\t%d\n", uData.UserId, uData.UserName, classTimes, workTimes, violationTimes)
		} else {
			fmt.Printf("%s\t%s\t%d\t\t%d\t\t%d\t%d\n", uData.UserId, uData.UserName, classTimes, workTimes, int(score/workTimes), violationTimes)
		}

	}
}
