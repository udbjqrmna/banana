package db

import (
	"fmt"
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
	lendConn         chan Connection  //借出的连接
	returnConn       chan Connection  //还回来的连接
	closeCh          chan bool        //关闭的通道
	running          bool             //是否运行状态
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
		return nil, errors.Mistake("核心池数量不能大于总数量。")
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
		lendConn:         make(chan Connection),
		returnConn:       make(chan Connection),
		closeCh:          make(chan bool),
		running:          true,
	}

	if err := pool.init(); err != nil {
		return nil, err
	}

	pools[name] = &pool
	return &pool, nil
}

func NewDefaultPool(connectionUrl string, maxCount uint8, coreCount uint8, createConnection CreateConnection) (*ConnPool, error) {
	return NewPool(Default, connectionUrl, maxCount, coreCount, createConnection)
}

//返回一个空闲的连接，此操作属于底层操作，在使用完之后必须进行归还
func (cp *ConnPool) GetConnect() Connection {
	//log.Trace().Msg("现在调用返回一个连接的方法")
	return <-cp.lendConn
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

	}
	cp.count = cp.coreCount

	//池开启协程
	go cp.run()
	return nil
}

//
func (cp *ConnPool) run() {
	log.Debug().Msg("数据库连接池开始运行")
	var newConn = cp.findIdleConnection()

	for cp.running {
		select {
		case rc := <-cp.returnConn:
			log.Trace().Msg(fmt.Sprintf("收到还回连接:%p", rc))
			rc.StopUse()
		case cp.lendConn <- newConn:
			log.Trace().Msg(fmt.Sprintf("已经将连接%p 给到请求者", newConn))
			newConn = cp.findIdleConnection()
		case <-cp.closeCh:
			cp.running = false
		}
	}

	log.Info().Msg("数据库连接池结束操作，正常返回")
	cp.destroy()
}

func (cp *ConnPool) findIdleConnection() Connection {
	//log.Debug().Msg("启动查找有效连接的过程。")
	i := cp.index

	for ; i < cp.count; i++ {
		if !cp.connections[i].Busy() {
			cp.index = i + 1
			cp.connections[i].StartUse()
			log.Trace().Msg(fmt.Sprintf("给出一个连接:%p，index:%v", cp.connections[i], i))
			return cp.connections[i]
		}
	}

	for i = 0; i < cp.index; i++ {
		if !cp.connections[i].Busy() {
			cp.index = i + 1
			cp.connections[i].StartUse()
			log.Trace().Msg(fmt.Sprintf("给出一个连接222:%p，index:%v", cp.connections[i], i))
			return cp.connections[i]
		}
	}

	//到这里就代表所有都忙
	//如果未达到最大池容量
	if cp.count < cp.maxCount {
		log.Debug().Msg("增加池容量")
		if conn, err := cp.createConnection(cp.connectUrl); err != nil {
			log.Error().Msg("新生成conn异常。错误：" + err.Error())
		} else {
			cp.connections[cp.count] = conn
			cp.count++
			cp.index = 0
			conn.StartUse()
			log.Trace().Msg(fmt.Sprintf("给出一个连接3:%p", cp.connections[cp.count-1]))
			return conn
		}
	}

	//log.Trace().Msg("现在没有连接了，在这里等着吧")
	//conn := <-cp.returnConn
	//conn.StopUse()
	//conn.StartUse()
	//return conn

	select {
	case conn := <-cp.returnConn:
		log.Trace().Msg(fmt.Sprintf("在方法内收到还回连接:%p", conn))
		conn.StopUse()
		conn.StartUse()
		return conn
	case <-cp.closeCh:
		cp.running = false
	}

	//没有了直接返回nil
	return nil
}

func (cp *ConnPool) destroy() {
	log.Trace().Msg("开始销毁动作")
	//只关闭发送者的通道
	close(cp.lendConn)

	//调用连接的关闭方法
	for i := uint8(0); i < cp.count; i++ {
		cp.connections[i].Close()
	}
}

func (cp *ConnPool) Close() {
	log.Trace().Msg("开始关闭操作")
	cp.closeCh <- true
}

func (cp *ConnPool) ReturnConnection(conn Connection) {
	cp.returnConn <- conn
}
