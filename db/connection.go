package db

//HandleConnected 定义回调函数的操作
type HandleConnected = func(c Connection) error

//Connection 连接对象。所有数据库操作都经过这里
type Connection interface {
	//创建一个新的连接
	New(connectionUrl string, h HandleConnected) (Connection, error)
	//返回当前连接的连接字符串
	GetConnectUrl() string
	//此连接使用完之后返回操作，在这里进行一些清理工作
	StopUse()
	//此连接被使用前必须调用此方法。
	StartUse()
	//返回当前连接是否正在被使用,true:正在被使用
	Busy() bool
	//关闭连接
	Close()
}

//HandleTableResult
//type HandleTableResult func(table Table) error
