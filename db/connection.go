package db

//HandleConnected 定义回调函数的操作
type HandleConnected = func(c Connection) error

//Connection 连接对象。所有数据库操作都经过这里
type Connection interface {
	//创建一个新的连接
	New(connectionUrl string, h HandleConnected) (Connection, error)
	GetConnectUrl() string
	//Query(sql string, h HandleTableResult) error
}

//HandleTableResult
//type HandleTableResult func(table Table) error
