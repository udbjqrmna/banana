package db

import (
	"fmt"
	"github.com/udbjqrmna/banana/errors"
	"time"
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
	idleShrinkTime   time.Duration    //连接池空闲时进行收缩的值
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
		return nil, errors.NotUnderstand("poolName")
	}
}

type CreateConnection = func(connectionUrl string) (Connection, error)

func NewPool(name, connectionUrl string, maxCount uint8, coreCount uint8, createConnection CreateConnection, idleShrinkTime uint8) (*ConnPool, error) {
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
		idleShrinkTime:   time.Duration(idleShrinkTime),
		running:          true,
	}

	go pool.run()

	pools[name] = &pool
	return &pool, nil
}

func NewDefaultPool(connectionUrl string, maxCount uint8, coreCount uint8, createConnection CreateConnection, idleShrinkTime uint8) (*ConnPool, error) {
	return NewPool(Default, connectionUrl, maxCount, coreCount, createConnection, idleShrinkTime)
}

//返回一个空闲的连接，此操作属于底层操作，在使用完之后必须进行归还
func (cp *ConnPool) GetConnect() Connection {
	//log.Trace().Msg("现在调用返回一个连接的方法")
	return <-cp.lendConn
}

//初始化池对象，此方法在连接池的协程操作内执行
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

	return nil
}

//shrink 到时间的收缩方法。此方法在连接池的协程操作内执行
func (cp *ConnPool) shrink() {
	//当已经在核心数量时，不再进行收缩
	if cp.count <= cp.coreCount {
		return
	}

	for i := cp.coreCount; i < cp.maxCount; i++ {
		conn := cp.connections[i]
		if conn != nil && !conn.Busy() {
			log.Trace().Msgf("找到一个要收缩的%v，进行收缩。", conn)
			conn.Close()

			cp.connections[i] = nil
			cp.count--
		}
	}
}

//这个方法将单独启动一个协程进行服务。在此方法下所有调用都必须在这个单独的协程里面实现
func (cp *ConnPool) run() {
	if err := cp.init(); err != nil {
		log.Error().Error(err).Msg("初始化出现错误，停止运行")
		return
	}

	log.Debug().Msg("数据库连接池服务开始运行")

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
		case <-time.After(time.Second * cp.idleShrinkTime):
			log.Trace().Msg("现在已经触发了空闲操作。")
			cp.shrink()
		}
	}

	log.Info().Msg("数据库连接池结束操作，正常返回")
	cp.destroy()
}

//findIdleConnection　查找一个空闲的连接，此方法在连接池的协程操作内执行
func (cp *ConnPool) findIdleConnection() Connection {
	i := cp.index

	for ; i < cp.count; i++ {
		if !cp.connections[i].Busy() {
			cp.index = i + 1
			cp.connections[i].StartUse()
			return cp.connections[i]
		}
	}

	for i = 0; i < cp.index; i++ {
		if !cp.connections[i].Busy() {
			cp.index = i + 1
			cp.connections[i].StartUse()
			return cp.connections[i]
		}
	}

	//到这里就代表所有都忙
	log.Debug().Msg("连接都在忙呀。到这里等一下吧。")

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

	select {
	case conn := <-cp.returnConn:
		log.Trace().Msg(fmt.Sprintf("在方法内收到还回连接:%p", conn))
		conn.StopUse()
		conn.StartUse()
		return conn
	case <-cp.closeCh:
		cp.running = false
		//关闭操作，所以直接返回nil
		return nil
	}
}

//destroy　在关闭后进行最后的清理
func (cp *ConnPool) destroy() {
	log.Trace().Msg("开始销毁动作")
	//只关闭发送者的通道
	close(cp.lendConn)

	//调用连接的关闭方法
	for i := uint8(0); i < cp.count; i++ {
		cp.connections[i].Close()
	}
}

//Close 结束连接池的关闭操作。此连接池关闭后将不能再继续操作。
func (cp *ConnPool) Close() {
	log.Trace().Msg("开始关闭操作")
	cp.closeCh <- true
}

//ReturnConnection 还回连接的方法
func (cp *ConnPool) ReturnConnection(conn Connection) {
	cp.returnConn <- conn
}
