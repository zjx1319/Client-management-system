package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
	"net"
	"src/config"
	"src/data"
	"src/tcp"
	"time"

	"github.com/kbinani/screenshot"
)

var VideoFlag bool = false

//发送屏幕截图
func sendScreenShot() {
	var msg data.Message
	msg.Type = data.ScreenShotResType
	var screenShotRes data.ScreenShotRes
	img, _ := screenshot.CaptureDisplay(0)
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	imgByte := buf.Bytes()
	screenShotRes.Img = base64.StdEncoding.EncodeToString(imgByte)
	dataByte, _ := json.Marshal(screenShotRes)
	msg.Data = string(dataByte)
	dataByte, _ = json.Marshal(msg)
	tcp.WritePkg(conn, dataByte)
}

//循环截图并比较屏幕内容
func screenCheckProcess() {
	var img1 *image.RGBA
	img2, _ := screenshot.CaptureDisplay(0)
	var msg data.Message
	msg.Type = data.ScreenReportType
	var screenReport data.ScreenReport
	var dataByte []byte
	for {
		time.Sleep(time.Minute) //延迟1分钟
		img1 = img2
		img2, _ = screenshot.CaptureDisplay(0)
		if imgCompare(img1, img2) < 1000 {
			screenReport.UnchangeTime++
		} else {
			screenReport.UnchangeTime = 0
		}
		//发送数据
		dataByte, _ = json.Marshal(screenReport)
		msg.Data = string(dataByte)
		dataByte, _ = json.Marshal(msg)
		tcp.WritePkg(conn, dataByte)
	}
}

//图像比较 结果大概1000以上说明有较明显变化
func imgCompare(img1, img2 *image.RGBA) (result int) {
	if len(img1.Pix) != len(img2.Pix) {
		//图像大小不相等 无法比较
		result = -1
		return
	}
	imgSize := len(img1.Pix)
	for i := 0; i < imgSize; i++ {
		diff := int(img1.Pix[i] - img2.Pix[i])
		result = result + diff*diff
	}
	result = result / imgSize
	return
}

func sendScreenVideo() {
	var msg data.Message
	msg.Type = data.ScreenVideoResType
	var screenVideoRes data.ScreenVideoRes
	for VideoFlag {
		img, _ := screenshot.CaptureDisplay(0)
		buf := new(bytes.Buffer)
		png.Encode(buf, img)
		imgByte := buf.Bytes()
		screenVideoRes.Img = base64.StdEncoding.EncodeToString(imgByte)
		dataByte, _ := json.Marshal(screenVideoRes)
		msg.Data = string(dataByte)
		dataByte, _ = json.Marshal(msg)
		//新建连接 主要是避免粘包
		videoConn, _ := net.Dial("tcp", config.Server)
		tcp.WritePkg(videoConn, dataByte)
		time.Sleep(time.Millisecond * 100)
	}
}
