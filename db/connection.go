package db

//HandleConnected 定义回调函数的操作
type HandleConnected = func(c Connection) error

//Connection 连接对象。所有数据库操作都经过这里
type Connection interface {
	//返回当前连接的连接字符串
	GetConnectUrl() string
	//此连接使用完之后返回操作，在这里进行一些清理工作
	StopUse()
	//此连接被使用前必须调用此方法。在此实现调用前的一些初始工作
	StartUse() error
	//返回当前连接是否正在被使用,true:正在被使用
	Busy() bool
	//关闭连接
	Close()
	//启动连接，如果已经连接，此方法将重新连接
	StartLink() error

	////启动一个事务，当前启动过事务无法再次启动事务
	//BegingTran()
	////提交事务
	//Commit()
	////回滚事务
	//Rollback()
	////执行一个查询
	//Query()
	//Execute()
	////复制操作
	//CopyInDb()
	//CopyOutDb()
	////返回当前连接的健康程度，0:未初始化，1:就绪，可发命令　5:在事务中　20:断开连接
	//Health()

	////预编译的操作
	//PreQuery()
	//PreExecute()
	////所有命令均为异步的操作
}

//HandleTableResult
//type HandleTableResult func(table Table) error
