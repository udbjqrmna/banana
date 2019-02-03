package postgresql

import (
	logD "github.com/udbjqrmna/banana/db/log"
	"time"
)

//定义此包的全局对象
var (
	//日志对象
	log = logD.Log()

	connectionTryIdleTime time.Duration = 3 //连接失败后尝试的间隔时间，单位：秒

	//定义发送消息
	ConnCloseMsg []byte
)

const (
	netProtocol = "tcp" //使用的网络协议
	defaultBuf  = 4096  //接收缓冲区默认大小

	//连接命令使用的一些常量
	connClose = 255 //连接关闭的命令值
)

func init() {
	initCommand()
}

//初始化命令格式
func initCommand() {
	cmd := make([]byte, 1)
	ConnCloseMsg = append(cmd, connClose)
}
