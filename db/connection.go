package db

type HandleConnected = func(Connection)

//Connection 连接对象。所有数据库操作都经过这里
type Connection interface {
	//创建一个新的连接
	create(connectionUrl string, h HandleConnected) (*Connection, error)
	//Query(sql string, h HandleTableResult) error
}

//HandleTableResult
//type HandleTableResult func(table Table) error
