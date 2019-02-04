package postgresql

import "github.com/udbjqrmna/banana/db/postgresql/protocol3"

const (
	authenticationKey = 'R'
	errorResponse     = 'E'
)

/*

这个文件就是为了连接的返回函数。
因为连接返回的处理较多，将此单独一个文件来处理，目的是更清晰
*/

func (c *Connection) handleResponse(responseMsg []byte) error {
	log.Trace().Bytes("rspMsg", responseMsg).Msg("开始处理服务端消息")
	switch responseMsg[0] {
	case authenticationKey:
		msg := protocol3.Authentication{}
		msg.Decode(responseMsg)
		switch msg.Style {
		case protocol3.StyleCleartext:
			_, err := c.pgConn.Write(protocol3.EncodePassword(c.connPara[password]))
			return err
		case protocol3.StyleMd5:
			digestedPassword := "md5" + hexMD5(hexMD5(c.connPara[password]+c.connPara[user])+string(msg.Salt))

			_, err := c.pgConn.Write(protocol3.EncodePassword(digestedPassword))
			return err
		}
	case errorResponse:
		log.Error().Msgf("收到一个错误。%s", responseMsg)
		msg := protocol3.ErrorResponse{}
		msg.Decode(responseMsg)
		switch msg.Code {
		//invalid_password or invalid_authorization_specification
		case "28P01", "28000":
			c.status = connError
			c.errMsg = msg.Body
			log.Error().String("code", msg.Code).String("err", c.errMsg).Msg("身份验证出现错误")
		}
	}

	//TODO 取第一个值，判断类型，并且返回对应的结果
	return nil
}
