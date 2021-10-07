package main

import (
	"src/data"

	"github.com/fatih/color"
)

func askQueRes(user data.User, queMes data.QueMes, seat string) {
	color.HiMagenta("[Question][%s]机位为%s的学生%s提出了一个问题：%s\n", user.UserId, seat, user.UserName, queMes.Content)
}
