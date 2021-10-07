package config

import "time"

const (
	//Student
	//Server string = "localhost:8989" //服务端ip和端口
	Server string = "192.168.177.1:8989" //服务端ip和端口

	//Teacher
	SerPoint         string        = "8989" //服务端端口
	RedisIp          string        = "127.0.0.1:6379"
	RedisMaxIdle     int           = 16                //最大空闲链接数
	RedisMaxActive   int           = 256               // 表示和数据库的最大链接数
	RedisIdleTimeout time.Duration = time.Second * 300 // 最大空闲时间
)

var Seat string = "seat" //机位
