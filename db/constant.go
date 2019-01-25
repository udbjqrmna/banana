package db

import (
	"errors"
	"github.com/udbjqrmna/onelog"
	"os"
)

const (
	//EmptyString 一个空字符串的值，不是nil
	EmptyString string = ""
	//定义一个默认的名称
	Default = "default"
	LogName = "DBLog"
)

//定义此包的全局对象
var (
	//NotAllowOperation 不允许的操作
	NotAllowOperation = errors.New("此操作不被允许")
	//日志对象
	log *onelog.Logger = nil
	//保存池对象的map
	pools = make(map[string]*ConnPool)
)

func Log() *onelog.Logger {
	return log
}

func init() {
	onelog.TimeFormat = "15:04:05.000000"
	//实例化日志对象
	if log = onelog.GetLog(LogName); log == nil {
		log = onelog.New(&onelog.Stdout{Writer: os.Stdout}, onelog.TraceLevel, &onelog.JsonPattern{}).AddRuntime(&onelog.CoroutineID{})
	}
}
