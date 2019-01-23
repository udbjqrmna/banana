package db

import (
	"github.com/udbjqrmna/banana/errors"
)

//ConnPool 连接池的接口，所有连接从这里获取
type ConnPool struct {
	Name             string           //连接池名称
	connectUrl       string           //连接字符串
	maxCount         uint8            //最大连接数，超过此数不再建立连接
	coreCount        uint8            //核心连接数，最小保持此连接数
	index            uint8            //当前池内对象的指针
	count            uint8            //当前的总可用数
	connections      []Connection     //连接slice
	createConnection CreateConnection //返回连接对象的操作
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

type CreateConnection = func(connectionUrl string) (Connection, error)

func NewPool(name, connectionUrl string, maxCount uint8, coreCount uint8, createConnection CreateConnection) (*ConnPool, error) {
	if coreCount > maxCount {
		return nil, errors.MisTake("核心池数量不能大于总数量。")
	}

	if _, ok := pools[name]; ok {
		return nil, errors.AlreadyExisting(name)
	}

	//构建连接池对象
	var pool = ConnPool{
		Name:             name,
		connectUrl:       connectionUrl,
		maxCount:         maxCount,
		coreCount:        coreCount,
		index:            0,
		count:            0,
		connections:      make([]Connection, maxCount),
		createConnection: createConnection,
	}

	if err := pool.init(); err != nil {
		return nil, err
	}

	pools[name] = &pool
	return &pool, errors.UnknownMistake("")
}

func NewDefaultPool(connectionUrl string, maxCount uint8, coreCount uint8, createConnection CreateConnection) (*ConnPool, error) {
	return NewPool(Default, connectionUrl, maxCount, coreCount, createConnection)
}

//返回一个空闲的连接，此操作属于底层操作，在使用完之后必须进行归还
func (cp *ConnPool) GetConnect() (*Connection, error) {
	//TODO 返回一个空闲的连接
	return nil, nil
}

//初始化池对象
func (cp *ConnPool) init() error {
	for i := uint8(0); i < cp.coreCount; i++ {
		log.Trace().Int("number", int(i)).Msg("开始进行连接的初始化")
		conn, err := cp.createConnection(cp.connectUrl)
		if err != nil {
			return err
		}

		cp.connections[i] = conn
		cp.count++
	}

	return nil
}
