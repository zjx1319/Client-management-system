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
	_, err = conn.Read(dataByte[:4])
	if err != nil {
		return
	}

	pkgLen := binary.BigEndian.Uint32(dataByte[0:4])
	//根据 pkgLen 读取消息内容
	n, err := conn.Read(dataByte[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}
	//反序列化
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
