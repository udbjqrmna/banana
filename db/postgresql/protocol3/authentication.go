package protocol3

import (
	"encoding/binary"
)

/*
此文件写关于处理登录部分的请求。
*/

//验证部分的一些常量
const (
	key1 = 8 //需要验证的第一个值
	key2 = 12

	ok                = 0 //登录成功
	kerberosV5        = 2 //以KerberosV5方式进行验证
	cleartextPassword = 3 //需要以一个明文方式提交验证
	mD5Password       = 5 //需要以MD5Password方式提交验证
	//sCMCredential     = 6 //以SCMCredential方式进行验证
	//gSS               = 7 //以GSS方式进行验证
	//sSP               = 9 //以SSP方式进行验证

	StyleCleartext = 0x10
	StyleMd5       = 0x20
	StyleOk        = 0x00
	StyleUnknown   = 0xEF
	StyleError     = 0xFF
)

type Authentication struct {
	backMessage
	Style int
	Salt  []byte
}

func (a *Authentication) Decode(data []byte) {
	key := binary.BigEndian.Uint32(data[1:])

	switch key {
	case key1:
		value := binary.BigEndian.Uint32(data[5:])
		switch value {
		case ok:
			log.Trace().Msg("收到登录验证成功消息")
			a.Style = StyleOk
			return
		case kerberosV5:
			log.Error().Msg("以KerberosV5方式进行验证。目前未知，未操作。直接返回")
			a.Style = StyleUnknown
			return
		case cleartextPassword:
			log.Trace().Msg("需要明文验证")
			a.Style = StyleCleartext
			return
		}
	case key2:
		value := binary.BigEndian.Uint32(data[5:])
		switch value {
		case mD5Password:
			a.Salt = data[9:]
			log.Trace().Msg("需要Md5验证")
			a.Style = StyleMd5
			return
		}
	default:
		log.Error().Msg("未处理的验证方式。未操作。直接返回")
		a.Style = StyleUnknown
		return
	}

	a.Style = StyleError
	return
}
