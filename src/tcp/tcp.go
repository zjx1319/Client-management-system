package tcp

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"src/data"
)

func ReadPkg(conn net.Conn) (mes data.Message, err error) {

	dataByte := make([]byte, 8096)
	//fmt.Println("读取客户端发送的数据...")
	_, err = conn.Read(dataByte[:4])
	if err != nil {
		return
	}
	//根据buf[:4] 转成一个 uint32类型

	pkgLen := binary.BigEndian.Uint32(dataByte[0:4])
	//根据 pkgLen 读取消息内容
	n, err := conn.Read(dataByte[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}
	//把pkgLen 反序列化成 -> message.Message
	// 技术就是一层窗户纸 &mes！！
	err = json.Unmarshal(dataByte[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarsha err=", err)
		return
	}
	return
}

func WritePkg(conn net.Conn, wdata []byte) (err error) {

	//先发送一个长度给对方
	pkgLen := uint32(len(wdata))
	var dataByte [4]byte
	binary.BigEndian.PutUint32(dataByte[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(dataByte[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	//发送data本身
	n, err = conn.Write(wdata)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}

/*
func ClientReadPackage(conn net.Conn) (msg message.Message, err error) {
	var data [8192]byte
	n, err := conn.Read(data[0:4])
	if n != 4 {
		fmt.Println("读取错误, 休息30秒，再读数据...")
		time.Sleep(time.Second * 30)
		return
	}
	//fmt.Println("read package:", buf[0:4])

	var packLen uint32
	packLen = binary.BigEndian.Uint32(data[0:4])

	//fmt.Printf("receive len:%d", packLen)
	n, err = conn.Read(data[0:packLen])
	if n != int(packLen) {
		//err = errors.New("read body failed")
		return
	}

	//fmt.Printf("receive data:%s\n", string(buf[0:packLen]))
	err = json.Unmarshal(data[0:packLen], &msg)
	if err != nil {
		fmt.Println("unmarshal failed, err:", err)
	}
	return
}

func ClientWritePkg(conn net.Conn,wdata []byte) (err error) {

	//先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(wdata))
	var data [4]byte
	binary.BigEndian.PutUint32(data[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(data[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	//发送data本身
	n, err = conn.Write(wdata)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}
*/
