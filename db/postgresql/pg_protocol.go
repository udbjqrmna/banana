package postgresql

import "encoding/binary"

type Message interface {
	//将消息转成本对象的值
	Decode(data []byte) error
	//Encode 将赋值好的对象转成协议理解的byte值。
	Encode() []byte
}

//抽象结构，此结构目的为后续代码提供一个抽象的实现。
type AbsMessage struct {
}

func (*AbsMessage) Decode(data []byte) error {
	log.Warn().Msg("此对象未实现此方法")
	return nil
}

func (*AbsMessage) Encode() []byte {
	log.Warn().Msg("此对象未实现此方法")
	return nil
}

//前端向后端发送的消息
type FrontMessage struct {
	AbsMessage
	buf []byte
}

//后端向前端发送的消息
type BackMessage struct {
	AbsMessage
}

//BothMessage 前后端都有的消息。
type BothMessage struct {
	AbsMessage
}

const (
	AuthenticationKey = 'R'
)

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

func AppendInt64(buf []byte, n int64) []byte {
	return AppendUint64(buf, uint64(n))
}

func SetInt32(buf []byte, n int32) {
	binary.BigEndian.PutUint32(buf, uint32(n))
}
