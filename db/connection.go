package db

//Connection 连接对象。所有数据库操作都经过这里
type Connection interface {
	Query(sql string, h HandleTableResult) error
}

//HandleTableResult
type HandleTableResult func(table Table) error
