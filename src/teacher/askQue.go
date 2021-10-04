package main

import (
	"fmt"
	"src/data"
)

func askQueRes(user data.User, queMes data.QueMes, seat string) {
	fmt.Printf("[Q][%s]机位为%s的学生%s提出了一个问题：%s\n", user.UserId, seat, user.UserName, queMes.Content)
}
