package postgresql

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"github.com/udbjqrmna/banana/db"
	"io"
	"strconv"
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
	//sCMCredential     = 6 //以SCMCredential方式进行验证
	//gSS               = 7 //以GSS方式进行验证
	//sSP               = 9 //以SSP方式进行验证
	mD5Password = 5 //需要以MD5Password方式提交验证
)

//handleAuthentication 当服务器返回登录验证的要求里，此方法被调用。
//在此方法内时行不同验证方式的处理
func (c *Connection) handleAuthentication(data []byte) error {
	key := binary.BigEndian.Uint32(data)

	switch key {
	case key1:
		value := binary.BigEndian.Uint32(data[4:])
		switch value {
		case ok:
			return nil
		case kerberosV5:
			log.Error().Msg("以KerberosV5方式进行验证。目前未知，未操作。直接返回")
			return db.Unknown("未进行操作的验证方式：kerberosV5")
		case cleartextPassword:
			return c.sendCleartextPassword()
		}
		log.Trace().Msg("收到登录成功消息。")
		return nil
	case key2:
		value := binary.BigEndian.Uint32(data[4:])
		switch value {
		case mD5Password:
			return c.sendMD5Password(data[9:])
		}
	default:
		log.Error().Msg("未处理的验证方式。未操作。直接返回")
		return db.Unknown("未处理的验证方式。未操作。直接返回")
	}

	return db.Unknown("未能理解值：" + strconv.FormatUint(uint64(key), 10))
}

//发送明文验证
func (c *Connection) sendCleartextPassword() error {
	_, err := c.conn.Write(passwdMessage([]byte(c.connPara[password])))
	return err
}

//发送MD5验证
func (c *Connection) sendMD5Password(md5 []byte) error {
	digestedPassword := "md5" + hexMD5(hexMD5(c.connPara[password]+c.connPara[user])+string(md5))

	_, err := c.conn.Write(passwdMessage([]byte(digestedPassword)))
	return err
}

func hexMD5(s string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, s)
	return hex.EncodeToString(hash.Sum(nil))
}
