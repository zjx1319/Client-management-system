package main

import (
	"encoding/json"
	"src/data"

	"github.com/garyburd/redigo/redis"
)

func userAdd(newUser data.User) {
	data, _ := json.Marshal(newUser)
	Rconn.Do("HSET", "userData", newUser.UserId, string(data))
}

func userLogin(loginMes data.LoginMes) (result string, userName string) {
	userName = ""
	res, err := redis.String(Rconn.Do("hget", "userData", loginMes.UserId))
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
