package protocol3

import (
	"bytes"
	"encoding/binary"
	logD "github.com/udbjqrmna/banana/db/log"
)

var log = logD.Log()

type Message interface {
	//将消息转成本对象的值
	Decode(data []byte)
	//Encode 将赋值好的对象转成协议理解的byte值。
	Encode() []byte
}

//前后端都有的消息直接实现这个抽象结构
type bothMessage struct {
}

func (*bothMessage) Decode(data []byte) {
	log.Warn().Msg("此对象未实现此方法")
	return
}

func (*bothMessage) Encode() []byte {
	log.Warn().Msg("此对象未实现此方法")
	return nil
}

//前端向后端发送的消息
type frontMessage struct {
	bothMessage
}

//后端向前端发送的消息
type backMessage struct {
	bothMessage
}

//BothMessage 前后端都有的消息。
type BothMessage struct {
	bothMessage
}

func AppendUint16(buf []byte, n uint16) []byte {
	wp := len(buf)
	buf = append(buf, 0, 0)
	binary.BigEndian.PutUint16(buf[wp:], n)
	return buf
}

func AppendUint32(buf []byte, n uint32) []byte {
	wp := len(buf)
	buf = append(buf, 0, 0, 0, 0)
	binary.BigEndian.PutUint32(buf[wp:], n)
	return buf
}

func AppendUint64(buf []byte, n uint64) []byte {
	wp := len(buf)
	buf = append(buf, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.BigEndian.PutUint64(buf[wp:], n)
	return buf
}

func AppendInt16(buf []byte, n int16) []byte {
	return AppendUint16(buf, uint16(n))
}

func AppendInt32(buf []byte, n int32) []byte {
	return AppendUint32(buf, uint32(n))
}

func AppendInt(buf []byte, n int) []byte {
	return AppendUint32(buf, uint32(n))
}

func AppendInt64(buf []byte, n int64) []byte {
	return AppendUint64(buf, uint64(n))
}

func SetInt32(buf []byte, n int32) {
	binary.BigEndian.PutUint32(buf, uint32(n))
}

func SetInt(buf []byte, n int) {
	binary.BigEndian.PutUint32(buf, uint32(n))
}

//此方法得到返回消息里的内容。
//同时返回内容和已经取完值的0所在位置
func takeContent(buf []byte, prefix []byte) (string, int) {
	ind := bytes.Index(buf, prefix)
	if ind == -1 {
		return "", -1
	}
	end := ind + len(prefix) + bytes.IndexByte(buf[ind+len(prefix):], 0)
	return string(buf[ind+len(prefix) : end]), end
}
