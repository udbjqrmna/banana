package log

import (
	"github.com/udbjqrmna/onelog"
	"os"
)

const (
	Name = "DBLog"
)

//定义此包的全局对象
var (
	log *onelog.Logger = nil //日志对象
)

func Log() *onelog.Logger {
	return log
}

func init() {
	onelog.TimeFormat = "01-02 15:04:05.000000"
	//实例化日志对象
	if log = onelog.GetLog(Name); log == nil {
		log = onelog.New(&onelog.Stdout{Writer: os.Stdout}, onelog.TraceLevel, &onelog.JsonPattern{}).AddRuntime(&onelog.CoroutineID{})
	}
}
