package db

import (
	"github.com/udbjqrmna/banana/errors"
	"github.com/udbjqrmna/onelog"
	"os"
)

var (
	//保存池对象的map
	pools = make(map[string]*ConnPool)
	//日志对象
	log *onelog.Logger = nil
)

func init() {
	if log = onelog.GetLog("DbLog"); log == nil {
		log = onelog.New(&onelog.Stdout{Writer: os.Stdout}, onelog.TraceLevel, &onelog.JsonPattern{}).AddRuntime(&onelog.Caller{})
	}
}

//ConnPool 连接池的接口，所有连接从这里获取
type ConnPool struct {
	Name        string
	connectUrl  string       //连接字符串
	maxCount    uint8        //最大连接数，超过此数不再建立连接
	coreCount   uint8        //核心连接数，最小保持此连接数
	index       uint8        //当前池内对象的指针
	count       uint8        //当前的总可用数
	connections []Connection //连接slice
}

//GetPool 得到默认的连接池对象，如果多个可使用命名对象
func GetPool() *ConnPool {
	return pools[Default]
}

//GetPoolByName 得到默认的连接池对象，如果多个可使用命名对象
func GetPoolByName(poolName string) (*ConnPool, error) {
	if pool, ok := pools[poolName]; ok {
		return pool, nil
	} else {
		return nil, errors.NotUnderstand("")
	}
}

func NewPool(connectionUrl string, maxCount uint8, coreCount uint8) (*ConnPool, error) {
	return nil, errors.UnknownMistake("")
}
